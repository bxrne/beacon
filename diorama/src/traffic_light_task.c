#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "traffic_light_task.h"
#include "cJSON.h"

extern QueueHandle_t event_queue;

typedef enum
{
  STATE_CAR_GREEN,
  STATE_CAR_YELLOW,
  STATE_CAR_RED,
  STATE_PED_GREEN
} traffic_state_t;

void log_telemetry(traffic_state_t state, bool pedestrian_waiting)
{
  char telemetry[128];
  const char *car_light = (state == STATE_CAR_GREEN) ? "GREEN" : (state == STATE_CAR_YELLOW) ? "YELLOW"
                                                             : (state == STATE_CAR_RED)      ? "RED"
                                                                                             : "RED";
  const char *ped_light = (state == STATE_PED_GREEN) ? "GREEN" : "RED";
  snprintf(telemetry, sizeof(telemetry), "{\"car_light\":\"%s\",\"pedestrian_light\":\"%s\",\"pedestrian_waiting\":%s}",
           car_light, ped_light, pedestrian_waiting ? "true" : "false");
  ESP_LOGI("TELEMETRY", "%s", telemetry);
}

void traffic_light_task(void *pvParameters)
{
  // Initial state setup
  traffic_state_t current_state = STATE_CAR_GREEN;
  bool pedestrian_waiting = false;

  // Set initial light states
  gpio_set_level(CAR_GREEN_PIN, 1);
  gpio_set_level(CAR_YELLOW_PIN, 0);
  gpio_set_level(CAR_RED_PIN, 0);
  gpio_set_level(PED_GREEN_PIN, 0);
  gpio_set_level(PED_RED_PIN, 1);

  ESP_LOGI("TRAFFIC_LIGHT_TASK", "Initial traffic light state set");

  while (1)
  {
    // Check for pedestrian button press
    event_t event;
    if (xQueueReceive(event_queue, &event, 0))
    {
      if (event == EVENT_BUTTON_PRESS)
      {
        ESP_LOGI("EVENT", "Pedestrian button pressed");
        pedestrian_waiting = true;
      }
    }

    switch (current_state)
    {
    case STATE_CAR_GREEN:
      gpio_set_level(CAR_GREEN_PIN, 1);
      gpio_set_level(CAR_YELLOW_PIN, 0);
      gpio_set_level(CAR_RED_PIN, 0);
      log_telemetry(current_state, pedestrian_waiting);

      if (pedestrian_waiting)
      {
        // Transition to yellow if pedestrian is waiting
        vTaskDelay(pdMS_TO_TICKS(500)); // Brief delay before transition
        current_state = STATE_CAR_YELLOW;
        ESP_LOGI("STATE", "Transitioning to yellow due to pedestrian");
      }
      else
      {
        vTaskDelay(pdMS_TO_TICKS(CAR_GREEN_DURATION));
        current_state = STATE_CAR_YELLOW;
      }
      break;

    case STATE_CAR_YELLOW:
      gpio_set_level(CAR_GREEN_PIN, 0);
      gpio_set_level(CAR_YELLOW_PIN, 1);
      ESP_LOGI("STATE", "Car light: Yellow");
      log_telemetry(current_state, pedestrian_waiting);
      vTaskDelay(pdMS_TO_TICKS(CAR_YELLOW_DURATION));
      current_state = STATE_CAR_RED;
      break;

    case STATE_CAR_RED:
      gpio_set_level(CAR_YELLOW_PIN, 0);
      gpio_set_level(CAR_RED_PIN, 1);
      ESP_LOGI("STATE", "Car light: Red");
      log_telemetry(current_state, pedestrian_waiting);

      if (pedestrian_waiting)
      {
        current_state = STATE_PED_GREEN;
      }
      else
      {
        vTaskDelay(pdMS_TO_TICKS(CAR_RED_DURATION));
        current_state = STATE_CAR_GREEN;
      }
      break;

    case STATE_PED_GREEN:
      gpio_set_level(PED_RED_PIN, 0);
      gpio_set_level(PED_GREEN_PIN, 1);
      ESP_LOGI("STATE", "Pedestrian light: Green");
      log_telemetry(current_state, pedestrian_waiting);
      vTaskDelay(pdMS_TO_TICKS(PED_GREEN_DURATION));

      // Transition back to red
      gpio_set_level(PED_GREEN_PIN, 0);
      gpio_set_level(PED_RED_PIN, 1);
      ESP_LOGI("STATE", "Pedestrian light: Red");

      pedestrian_waiting = false;
      current_state = STATE_CAR_GREEN;
      break;
    }
  }
}