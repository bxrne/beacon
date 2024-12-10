#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "freertos/timers.h"
#include "freertos/semphr.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "button_task.h"
#include "traffic_light_task.h"
#include "cJSON.h"
#include <string.h>     //for handling strings
#include "esp_system.h" //esp_init funtions esp_err_t
#include "esp_log.h"    //for showing logs
#include "nvs_flash.h"  //non volatile storage
#include "lwip/err.h"   //light weight ip packets error handling
#include "lwip/sys.h"   //system applications for light weight ip apps

QueueHandle_t event_queue;
SemaphoreHandle_t button_semaphore;
TaskHandle_t button_task_handle;
TaskHandle_t traffic_light_task_handle;

void app_main(void)
{
    esp_err_t ret = nvs_flash_init();
    if (ret == ESP_ERR_NVS_NO_FREE_PAGES || ret == ESP_ERR_NVS_NEW_VERSION_FOUND)
    {
        ESP_ERROR_CHECK(nvs_flash_erase());
        ret = nvs_flash_init();
    }
    ESP_ERROR_CHECK(ret);

    ESP_LOGI("APP_MAIN", "Starting application");

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
    ESP_LOGI("APP_MAIN", "Button pin configured with pull-up resistor");

    // Create event queue
    event_queue = xQueueCreate(EVENT_QUEUE_SIZE, sizeof(event_t));
    if (event_queue == NULL)
    {
        ESP_LOGE("APP_MAIN", "Failed to create event queue");
        return;
    }

    // Create button semaphore
    button_semaphore = xSemaphoreCreateBinary();
    if (button_semaphore == NULL)
    {
        ESP_LOGE("APP_MAIN", "Failed to create button semaphore");
        return;
    }

    // Create tasks with adjusted stack sizes and priorities
    if (xTaskCreate(traffic_light_task, "Traffic Light Task", TRAFFIC_LIGHT_TASK_STACK_SIZE, NULL, TRAFFIC_LIGHT_TASK_PRIORITY, &traffic_light_task_handle) != pdPASS)
    {
        ESP_LOGE("APP_MAIN", "Failed to create Traffic Light Task");
    }

    if (xTaskCreate(button_task, "Button Task", BUTTON_TASK_STACK_SIZE, NULL, BUTTON_TASK_PRIORITY, &button_task_handle) != pdPASS)
    {
        ESP_LOGE("APP_MAIN", "Failed to create Button Task");
    }

    // Initial traffic light state
    gpio_set_level(CAR_GREEN_PIN, 1);
    gpio_set_level(CAR_YELLOW_PIN, 0);
    gpio_set_level(CAR_RED_PIN, 0);
    gpio_set_level(PED_GREEN_PIN, 0);
    gpio_set_level(PED_RED_PIN, 1);
    ESP_LOGI("APP_MAIN", "Initial traffic light state set");

    // Delete this task if no longer needed
    vTaskDelete(NULL);
}
