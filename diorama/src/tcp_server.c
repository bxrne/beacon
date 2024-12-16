#include "tcp_server.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "lwip/sockets.h"
#include "esp_log.h"
#include <string.h> // For strlen and strstr

#define TAG "TCP_SERVER"

// Custom Protocol Design:
// - Start Byte: 0x02
// - Length Byte: Specifies the length of the payload
// - Payload: Actual data
// - End Byte: 0x03 (optional)

void tcp_server_task(void *pvParameters)
{
  // Create socket, bind, listen, and accept connections
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
  dest_addr.sin_port = htons(80);
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

    char rx_buffer[128];
    int len = recv(sock, rx_buffer, sizeof(rx_buffer) - 1, 0);
    if (len < 0)
    {
      ESP_LOGE(TAG, "recv failed: errno %d", errno);
      close(sock);
      continue;
    }

    // Check if the request is a GET request
    if (strstr(rx_buffer, "GET / ") != NULL)
    {
      // Prepare the response in the custom format
      char response[128];
      const char *payload = "Hello, World!";
      uint8_t payload_length = strlen(payload);

      response[0] = 0x02;           // Start Byte
      response[1] = payload_length; // Length Byte
      memcpy(&response[2], payload, payload_length);
      response[2 + payload_length] = 0x03; // End Byte (optional)

      send(sock, response, 3 + payload_length, 0); // Send the response
    }
    else
    {
      // Parse the custom protocol
      if (len >= 2 && rx_buffer[0] == 0x02) // Check for Start Byte
      {
        uint8_t payload_length = rx_buffer[1]; // Length Byte

        if (len >= 2 + payload_length)
        {
          // Extract the payload
          char payload[128];
          memcpy(payload, &rx_buffer[2], payload_length);
          payload[payload_length] = '\0';

          // Process the payload
          // For example, log the received payload
          ESP_LOGI(TAG, "Received payload: %s", payload);

          // ...additional payload processing...

          // Send acknowledgment in the specified format
          char response[128];
          response[0] = 0x02; // Start Byte
          response[1] = 3;    // Length Byte (length of "ACK")
          memcpy(&response[2], "ACK", 3);
          response[5] = 0x03; // End Byte (optional)

          send(sock, response, 6, 0); // Send the response
        }
        else
        {
          // Handle incorrect length
          ESP_LOGE(TAG, "Payload length mismatch");
        }
      }
      else
      {
        // Handle invalid start byte
        ESP_LOGE(TAG, "Invalid start byte");
      }
    }

    close(sock);
  }

  close(listen_sock);
  vTaskDelete(NULL);
}