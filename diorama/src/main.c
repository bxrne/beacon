#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>

#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "freertos/timers.h"
#include "freertos/semphr.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"

#define TAG "TRAFFIC_LIGHT"

QueueHandle_t pedestrianRequestQueue;
SemaphoreHandle_t xPedestrianSemaphore;

void IRAM_ATTR button_isr_handler(void *arg)
{
    uint32_t button_pressed = 1;
    xQueueSendFromISR(pedestrianRequestQueue, &button_pressed, NULL);
}

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

void CarLightTask(void *pvParameters)
{
    while (true)
    {
        // Car Green Light
        gpio_set_level(CAR_GREEN_PIN, 1);
        ESP_LOGI(TAG, "Car light: GREEN");
        vTaskDelay(pdMS_TO_TICKS(CAR_GREEN_DURATION));
        gpio_set_level(CAR_GREEN_PIN, 0);

        // Car Yellow Light
        gpio_set_level(CAR_YELLOW_PIN, 1);
        ESP_LOGI(TAG, "Car light: YELLOW");
        vTaskDelay(pdMS_TO_TICKS(CAR_YELLOW_DURATION));
        gpio_set_level(CAR_YELLOW_PIN, 0);

        // Car Red Light
        gpio_set_level(CAR_RED_PIN, 1);
        ESP_LOGI(TAG, "Car light: RED");

        if (xSemaphoreTake(xPedestrianSemaphore, 0) == pdTRUE)
        {
            ESP_LOGI(TAG, "Pedestrian light: GREEN");
            gpio_set_level(PED_RED_PIN, 0);
            gpio_set_level(PED_GREEN_PIN, 1);
            vTaskDelay(pdMS_TO_TICKS(PED_GREEN_DURATION));
            gpio_set_level(PED_GREEN_PIN, 0);
            gpio_set_level(PED_RED_PIN, 1);
            ESP_LOGI(TAG, "Pedestrian light: RED");
        }
        else
        {
            vTaskDelay(pdMS_TO_TICKS(CAR_RED_DURATION));
        }

        gpio_set_level(CAR_RED_PIN, 0);
    }
}

void PedestrianLightTask(void *pvParameters)
{
    uint32_t button_pressed;
    while (true)
    {
        if (xQueueReceive(pedestrianRequestQueue, &button_pressed, portMAX_DELAY))
        {
            ESP_LOGI(TAG, "Pedestrian button pressed");
            // Signal the CarLightTask
            xSemaphoreGive(xPedestrianSemaphore);
        }
    }
}

void app_main(void)
{
    init_gpio();

    pedestrianRequestQueue = xQueueCreate(10, sizeof(uint32_t));
    xPedestrianSemaphore = xSemaphoreCreateBinary();

    xTaskCreate(CarLightTask, "CarLightTask", 2048, NULL, 1, NULL);
    xTaskCreate(PedestrianLightTask, "PedestrianLightTask", 2048, NULL, 1, NULL);
}