#include <freertos/FreeRTOS.h>
#include "traffic_light.h"
#include "freertos/task.h"
#include "freertos/semphr.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "globals.h"

#define TAG "TRAFFIC_LIGHT"

void TrafficLightTask(void *pvParameters)
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

    // Handle pedestrian crossing (ISR gives this)
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
