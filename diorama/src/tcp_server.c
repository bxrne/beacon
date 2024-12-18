#include "tcp_server.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "lwip/sockets.h"
#include "esp_log.h"
#include <string.h>
#include "metrics.h"
#include "esp_system.h"
#include "driver/gpio.h"
#include "config.h"
#define TAG "TCP_SERVER"
#define START_BYTE 0x02
#define END_BYTE 0x03

// INFO: Custom Protocol Design:
// - Start Byte: 0x02
// - Length Byte: Specifies the length of the payload
// - Payload: Actual data
// - End Byte: 0x03 (optional as length exists but good for validation)
// key: value, key: value, key: value => (car_light: green, ped_light: red, recorded_at: UTC)

void TCPServerTask(void *pvParameters)
{
  int listen_sock = socket(AF_INET, SOCK_STREAM, IPPROTO_IP);
  if (listen_sock < 0)
  {
    ESP_LOGE(TAG, "Unable to create socket: errno %d", errno);
    vTaskDelete(NULL);
    return;
  }

  struct sockaddr_in dest_addr;
  dest_addr.sin_addr.s_addr = htonl(INADDR_ANY);
  dest_addr.sin_family = AF_INET;
  dest_addr.sin_port = htons(TCP_PORT);
  if (bind(listen_sock, (struct sockaddr *)&dest_addr, sizeof(dest_addr)) < 0)
  {
    ESP_LOGE(TAG, "Socket unable to bind: errno %d", errno);
    close(listen_sock);
    vTaskDelete(NULL);
    return;
  }

  if (listen(listen_sock, 1) < 0)
  {
    ESP_LOGE(TAG, "Error occurred during listen: errno %d", errno);
    close(listen_sock);
    vTaskDelete(NULL);
    return;
  }

  while (1)
  {
    struct sockaddr_in6 source_addr;
    socklen_t addr_len = sizeof(source_addr);
    int sock = accept(listen_sock, (struct sockaddr *)&source_addr, &addr_len);
    if (sock < 0)
    {
      ESP_LOGE(TAG, "Unable to accept connection: errno %d", errno);
      break;
    }

    char rx_buffer[1024];
    int total_len = 0;
    bool headers_received = false;
    int content_length = 0;

    while (1) // may come in chunks so loop
    {
      int len = recv(sock, rx_buffer + total_len, sizeof(rx_buffer) - total_len - 1, 0);
      if (len < 0)
      {
        ESP_LOGE(TAG, "recv failed: errno %d", errno);
        break;
      }
      else if (len == 0)
      {
        // Connection closed
        ESP_LOGI(TAG, "Connection closed");
        break;
      }
      else
      {
        total_len += len;
        rx_buffer[total_len] = '\0';

        if (!headers_received)
        {
          char *headers_end = strstr(rx_buffer, "\r\n\r\n");
          if (headers_end)
          {
            headers_received = true;
            char *cl_ptr = strstr(rx_buffer, "Content-Length:");
            if (cl_ptr)
            {
              cl_ptr += strlen("Content-Length:");
              while (*cl_ptr == ' ')
                cl_ptr++; // Skip spaces
              content_length = atoi(cl_ptr);
            }

            int headers_size = headers_end + 4 - rx_buffer;
            int body_len = total_len - headers_size;

            // If we've already received the whole body, exit loop
            if (body_len >= content_length)
            {
              break;
            }
          }
        }
        else
        {
          // If headers already received, check if we've got the entire body
          char *headers_end = strstr(rx_buffer, "\r\n\r\n");
          int headers_size = headers_end + 4 - rx_buffer;
          int body_len = total_len - headers_size;

          if (body_len >= content_length)
          {
            break;
          }
        }

        // If buffer is full and we haven't received all data, break
        if (total_len >= sizeof(rx_buffer) - 1)
        {
          ESP_LOGE(TAG, "Request too large");
          break;
        }
      }
    }

    if (total_len <= 0)
    {
      close(sock);
      continue;
    }

    ESP_LOGI(TAG, "Received request:\n%s", rx_buffer);

    // INFO: Check if the request is a POST request to /cmd
    if (strstr(rx_buffer, "POST /cmd") != NULL)
    {
      // Find end of headers
      char *body = strstr(rx_buffer, "\r\n\r\n");
      if (body != NULL)
      {
        body += 4; // Move past \r\n\r\n"

        int headers_size = body - rx_buffer;
        int body_len = total_len - headers_size;

        if (body_len < content_length)
        {
          ESP_LOGE(TAG, "Incomplete body received");
        }

        ESP_LOGI(TAG, "Received command: %.*s", content_length, body);

        char command[128];
        strncpy(command, body, content_length);
        command[content_length] = '\0';

        if (strcmp(command, "reboot") == 0)
        {
          // Send acknowledgment to client
          const char *response =
              "HTTP/1.1 200 OK\r\n"
              "Content-Type: text/plain\r\n"
              "Content-Length: 20\r\n"
              "\r\n"
              "Rebooting device...\n";
          send(sock, response, strlen(response), 0);
          ESP_LOGI(TAG, "Rebooting device on command.");
          close(sock);
          esp_restart();
        }
        else
        {
          // Handle unknown commands
          const char *response =
              "HTTP/1.1 400 Bad Request\r\n"
              "Content-Type: text/plain\r\n"
              "Content-Length: 15\r\n"
              "\r\n"
              "Unknown command\n";
          send(sock, response, strlen(response), 0);
          ESP_LOGI(TAG, "Unknown command received: %s", command);
        }
      }
      else
      {
        // Handle missing body
        const char *response =
            "HTTP/1.1 400 Bad Request\r\n"
            "Content-Type: text/plain\r\n"
            "Content-Length: 19\r\n"
            "\r\n"
            "No command received\n";
        send(sock, response, strlen(response), 0);
        ESP_LOGI(TAG, "No command in request.");
      }
    }
    else if (strstr(rx_buffer, "GET /metric") != NULL)
    {
      char payload[256];

      // Get the current light states and time
      LightColor car_light_state = get_recent_car_light_state();
      LightColor ped_light_state = get_recent_ped_light_state();
      const char *car_light_str = light_color_to_string(car_light_state);
      const char *ped_light_str = light_color_to_string(ped_light_state);
      char time_str[64];
      get_current_time_utc(time_str, sizeof(time_str));

      snprintf(payload, sizeof(payload), "car_light: %s, ped_light: %s, recorded_at: %s",
               car_light_str, ped_light_str, time_str);

      // Create the custom protocol response
      char custom_response[512];
      int payload_length = strlen(payload);
      custom_response[0] = START_BYTE;
      custom_response[1] = payload_length;
      memcpy(&custom_response[2], payload, payload_length);
      custom_response[2 + payload_length] = END_BYTE;

      send(sock, custom_response, 3 + payload_length, 0);
    }
    else if (strstr(rx_buffer, "GET / ") != NULL) // Sanity check
    {
      char response[128];
      const char *payload = "Hello, World!";
      uint8_t payload_length = strlen(payload);

      response[0] = START_BYTE;     // Start Byte
      response[1] = payload_length; // Length Byte
      memcpy(&response[2], payload, payload_length);
      response[2 + payload_length] = END_BYTE; // End Byte
      ESP_LOGI(TAG, "Sending payload: %s", payload);
      send(sock, response, 3 + payload_length, 0);
    }
    else
    {
      const char *response =
          "HTTP/1.1 404 Not Found\r\n"
          "Content-Type: text/plain\r\n"
          "Content-Length: 10\r\n"
          "\r\n"
          "Not Found\n";
      send(sock, response, strlen(response), 0);
    }

    close(sock);
  }

  close(listen_sock);
  vTaskDelete(NULL);
}