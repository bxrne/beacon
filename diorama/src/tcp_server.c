
#include "tcp_server.h"
#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "lwip/sockets.h"
#include "esp_log.h"
#include <string.h> // For strlen and strstr

#define TAG "TCP_SERVER"

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
    uint addr_len = sizeof(source_addr);
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

    rx_buffer[len] = 0;
    if (strstr(rx_buffer, "GET / ") != NULL)
    {
      const char *response = "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nok";
      send(sock, response, strlen(response), 0);
    }

    close(sock);
  }

  close(listen_sock);
  vTaskDelete(NULL);
}