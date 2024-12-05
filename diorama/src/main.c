#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "freertos/timers.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "button_task.h"
#include "traffic_light_task.h"
#include "cJSON.h"
#include <string.h>     //for handling strings
#include "esp_system.h" //esp_init funtions esp_err_t
#include "esp_wifi.h"   //esp_wifi_init functions and wifi operations
#include "esp_log.h"    //for showing logs
#include "esp_event.h"  //for wifi event
#include "nvs_flash.h"  //non volatile storage
#include "lwip/err.h"   //light weight ip packets error handling
#include "lwip/sys.h"   //system applications for light weight ip apps
#include "freertos/event_groups.h"
#include "esp_wifi.h"
#include "esp_log.h"
#include "nvs_flash.h"
#include "telemetry_task.h"
#include "http_client.h"
#include "wifi_task.h" // Include wifi_task.h

QueueHandle_t event_queue;
TaskHandle_t button_task_handle;
TaskHandle_t traffic_light_task_handle;
extern TaskHandle_t wifi_task_handle; // Added WiFi task handle

/* The event group allows multiple bits for each event, but we only care about two events:
 * - we are connected to the AP with an IP
 * - we failed to connect after the maximum amount of retries */
// Remove the following lines
// const int WIFI_CONNECTED_BIT = BIT0;
// const int WIFI_FAIL_BIT = BIT1;
static EventGroupHandle_t s_wifi_event_group;

// Modify event_handler to handle disconnection reasons
static void event_handler(void *arg, esp_event_base_t event_base, int32_t event_id, void *event_data)
{
    if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_START)
    {
        esp_wifi_connect();
    }
    else if (event_base == WIFI_EVENT && event_id == WIFI_EVENT_STA_DISCONNECTED)
    {
        wifi_event_sta_disconnected_t *disconn = (wifi_event_sta_disconnected_t *)event_data;
        ESP_LOGI("WIFI_EVENT", "Disconnected (Reason: %d)", disconn->reason);
        esp_wifi_connect(); // Keep trying to reconnect
    }
    else if (event_base == IP_EVENT && event_id == IP_EVENT_STA_GOT_IP)
    {
        ip_event_got_ip_t *event = (ip_event_got_ip_t *)event_data;
        ESP_LOGI("WIFI_EVENT", "Got IP: " IPSTR, IP2STR(&event->ip_info.ip));
    }
    else
    {
        ESP_LOGE("WIFI_EVENT", "Unexpected event: %ld", event_id); // Use %ld for int32_t
    }
}

// Remove the following function
// void wifi_init_sta(void)
// {
//     s_wifi_event_group = xEventGroupCreate();

//     esp_netif_init();
//     esp_event_loop_create_default();
//     esp_netif_create_default_wifi_sta();

//     wifi_init_config_t cfg = WIFI_INIT_CONFIG_DEFAULT();
//     esp_wifi_init(&cfg);

//     esp_event_handler_register(WIFI_EVENT, ESP_EVENT_ANY_ID, &event_handler, NULL);
//     esp_event_handler_register(IP_EVENT, IP_EVENT_STA_GOT_IP, &event_handler, NULL);

//     wifi_config_t wifi_config = {
//         .sta = {
//             .ssid = W_SSID,
//             .password = W_PASS,
//         },
//     };

//     esp_wifi_set_mode(WIFI_MODE_STA);
//     esp_wifi_set_config(ESP_IF_WIFI_STA, &wifi_config);
//     esp_wifi_start();

//     ESP_LOGI("WIFI_INIT", "WiFi initialization finished.");
// }

// Modify control_loop to run regardless of WiFi status
static void control_loop(void *pvParameters)
{
    ESP_LOGI("CONTROL_LOOP", "Starting control loop tasks");

    // Initialize GPIOs
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

    gpio_set_pull_mode(PED_BUTTON_PIN, GPIO_PULLUP_ONLY);
    ESP_LOGI("CONTROL_LOOP", "Button pin configured with pull-up resistor");

    // Create event queue
    event_queue = xQueueCreate(10, sizeof(event_t));
    if (event_queue == NULL)
    {
        ESP_LOGE("CONTROL_LOOP", "Failed to create event queue");
        return;
    }

    // Create tasks
    if (xTaskCreate(traffic_light_task, "Traffic Light Task", 8192, NULL, 1, &traffic_light_task_handle) != pdPASS)
    {
        ESP_LOGE("CONTROL_LOOP", "Failed to create Traffic Light Task");
    }

    if (xTaskCreate(button_task, "Button Task", 4096, NULL, 1, &button_task_handle) != pdPASS)
    {
        ESP_LOGE("CONTROL_LOOP", "Failed to create Button Task");
    }

    if (xTaskCreate(telemetry_task, "Telemetry Task", 8192, NULL, 1, NULL) != pdPASS)
    {
        ESP_LOGE("CONTROL_LOOP", "Failed to create Telemetry Task");
    }

    // Initial traffic light state
    gpio_set_level(CAR_GREEN_PIN, 1);
    gpio_set_level(CAR_YELLOW_PIN, 0);
    gpio_set_level(CAR_RED_PIN, 0);
    gpio_set_level(PED_GREEN_PIN, 0);
    gpio_set_level(PED_RED_PIN, 1);
    ESP_LOGI("CONTROL_LOOP", "Initial traffic light state set");

    while (1)
    {
        // Main control loop tasks
        vTaskDelay(pdMS_TO_TICKS(1000));
    }
}

void app_main(void)
{
    esp_err_t ret = nvs_flash_init();
    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND)
    {
        ESP_ERROR_CHECK(nvs_flash_erase());
        ret = nvs_flash_init();
    }
    ESP_ERROR_CHECK(ret);

    ESP_LOGI("APP_MAIN", "ESP_WIFI_MODE_STA");

    // Start WiFi task
    xTaskCreate(wifi_task, "WiFi Task", 4096, NULL, 1, &wifi_task_handle);

    // Start control loop
    xTaskCreate(control_loop, "Control Loop", 8192, NULL, 1, NULL);
}
