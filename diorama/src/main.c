#include <freertos/FreeRTOS.h>
#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

#include "freertos/task.h"
#include "freertos/queue.h"
#include "freertos/timers.h"
#include "freertos/semphr.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "car_light.h"
#include "ped_light.h"
#include "ped_request.h"
#include "globals.h"
#include "esp_wifi.h"
#include "esp_event.h"
#include "nvs_flash.h"
#include "lwip/sockets.h"
#include "esp_netif.h"

#define TAG "MAIN"

QueueHandle_t pedestrianRequestQueue;
SemaphoreHandle_t xPedestrianSemaphore;

void init_gpio(void)
{
    ESP_LOGI(TAG, "Initializing GPIOs");

    esp_rom_gpio_pad_select_gpio(CAR_GREEN_PIN);
    esp_rom_gpio_pad_select_gpio(CAR_YELLOW_PIN);
    esp_rom_gpio_pad_select_gpio(CAR_RED_PIN);
    esp_rom_gpio_pad_select_gpio(PED_GREEN_PIN);
    esp_rom_gpio_pad_select_gpio(PED_RED_PIN);
    esp_rom_gpio_pad_select_gpio(PED_BUTTON_PIN);

    gpio_set_direction(CAR_GREEN_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(CAR_YELLOW_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(CAR_RED_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(PED_GREEN_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(PED_RED_PIN, GPIO_MODE_OUTPUT);
    gpio_set_direction(PED_BUTTON_PIN, GPIO_MODE_INPUT);

    gpio_set_intr_type(PED_BUTTON_PIN, GPIO_INTR_POSEDGE);
    gpio_install_isr_service(0);
    gpio_isr_handler_add(PED_BUTTON_PIN, button_isr_handler, NULL);

    gpio_set_level(PED_RED_PIN, 1); // Ensure pedestrian light is red by default
}

void wifi_init(void)
{
    // Initialize NVS
    esp_err_t ret = nvs_flash_init();
    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND)
    {
        ESP_ERROR_CHECK(nvs_flash_erase());
        ret = nvs_flash_init();
    }
    ESP_ERROR_CHECK(ret);

    // Initialize the underlying TCP/IP stack
    ESP_ERROR_CHECK(esp_netif_init());

    // Create the default event loop
    ESP_ERROR_CHECK(esp_event_loop_create_default());

    // Initialize Wi-Fi including netif with default config
    esp_netif_create_default_wifi_sta();
    wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();
    ESP_ERROR_CHECK(esp_wifi_init(&cfg));

    // Set Wi-Fi configuration
    wifi_config_t wifi_config = {
        .sta = {
            .ssid = "SKYWYBGL",
            .password = "iZzsVpcMbEq4",
        },
    };
    ESP_ERROR_CHECK(esp_wifi_set_mode(WIFI_MODE_STA));
    ESP_ERROR_CHECK(esp_wifi_set_config(ESP_IF_WIFI_STA, &wifi_config));
    ESP_ERROR_CHECK(esp_wifi_start());

    // Connect to the Wi-Fi network
    ESP_ERROR_CHECK(esp_wifi_connect());
}

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

void app_main(void)
{
    ESP_LOGI(TAG, "Starting application");

    pedestrianRequestQueue = xQueueCreate(10, sizeof(uint32_t));
    xPedestrianSemaphore = xSemaphoreCreateBinary();

    init_gpio();
    init_ped_request();

    xTaskCreate(CarLightTask, "CarLightTask", 2048, NULL, 1, NULL);
    xTaskCreate(PedestrianLightTask, "PedestrianLightTask", 2048, NULL, 1, NULL);

    wifi_init();
    xTaskCreate(tcp_server_task, "tcp_server", 4096, NULL, 5, NULL);
}