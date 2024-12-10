#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>

#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "driver/gpio.h"
#include "esp_log.h"
#include "config.h"
#include "traffic_light_task.h"
#include "cJSON.h"
#include "esp_http_client.h"
#include "esp_netif.h" // For network interface functions

extern QueueHandle_t event_queue;

typedef enum
{
  STATE_CAR_GREEN,
  STATE_CAR_YELLOW,
  STATE_CAR_RED,
  STATE_PED_GREEN
} traffic_state_t;

static void send_metric(const char *state)
{
  // Create JSON payload
  cJSON *root = cJSON_CreateObject();
  cJSON *metrics_array = cJSON_AddArrayToObject(root, "metrics");
  cJSON *metric = cJSON_CreateObject();
  cJSON_AddItemToArray(metrics_array, metric);
  cJSON_AddStringToObject(metric, "type", "traffic_light");
  cJSON_AddStringToObject(metric, "unit", "color");
  cJSON_AddStringToObject(metric, "value", state);

  char *post_data = cJSON_PrintUnformatted(root);
  cJSON_Delete(root);

  esp_http_client_config_t config = {
      .url = API_URL,
  };
  esp_http_client_handle_t client = esp_http_client_init(&config);

  // Set headers
  esp_http_client_set_header(client, "Content-Type", "application/json");
  esp_http_client_set_header(client, "X-DeviceID", DEVICE_ID);

  // Set post data
  esp_http_client_set_method(client, HTTP_METHOD_POST);
  esp_http_client_set_post_field(client, post_data, strlen(post_data));

  // Perform the request
  esp_err_t err = esp_http_client_perform(client);
  if (err == ESP_OK)
  {
    int status_code = esp_http_client_get_status_code(client);
    ESP_LOGI("HTTP_CLIENT", "HTTP POST Status = %d", status_code);
  }
  else
  {
    ESP_LOGE("HTTP_CLIENT", "HTTP POST request failed: %s", esp_err_to_name(err));
  }

  // Clean up
  esp_http_client_cleanup(client);
  free(post_data);
}

void traffic_light_task(void *pvParameters)
{
  // Initial state setup
  traffic_state_t current_state = STATE_CAR_GREEN;
  bool pedestrian_button = false;

  ESP_LOGI("TRAFFIC_LIGHT_TASK", "Traffic light task started");

  while (1)
  {
    event_t event;
    // Check for button press event
    if (xQueueReceive(event_queue, &event, pdMS_TO_TICKS(100)))
    {
      if (event == EVENT_BUTTON_PRESS)
      {
        ESP_LOGI("TRAFFIC_LIGHT_TASK", "Pedestrian button pressed");
        pedestrian_button = true;
      }
    }

    switch (current_state)
    {
    case STATE_CAR_GREEN:
      gpio_set_level(CAR_GREEN_PIN, 1);
      gpio_set_level(CAR_YELLOW_PIN, 0);
      gpio_set_level(CAR_RED_PIN, 0);
      gpio_set_level(PED_GREEN_PIN, 0);
      gpio_set_level(PED_RED_PIN, 1);
      send_metric("pedestrian_red");
      ESP_LOGI("STATE", "Car light: Green");
      send_metric("car_green");

      if (pedestrian_button)
      {
        vTaskDelay(pdMS_TO_TICKS(500));
        current_state = STATE_CAR_YELLOW;
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
      send_metric("car_yellow");
      vTaskDelay(pdMS_TO_TICKS(CAR_YELLOW_DURATION));
      current_state = STATE_CAR_RED;
      break;

    case STATE_CAR_RED:
      gpio_set_level(CAR_YELLOW_PIN, 0);
      gpio_set_level(CAR_RED_PIN, 1);
      ESP_LOGI("STATE", "Car light: Red");
      send_metric("car_red");

      if (pedestrian_button)
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
      send_metric("pedestrian_green");
      vTaskDelay(pdMS_TO_TICKS(PED_GREEN_DURATION));

      gpio_set_level(PED_GREEN_PIN, 0);
      gpio_set_level(PED_RED_PIN, 1);
      send_metric("pedestrian_red");
      pedestrian_button = false;
      current_state = STATE_CAR_GREEN;
      break;
    }
  }
}