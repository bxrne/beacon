#include "freertos/FreeRTOS.h"
#include "freertos/task.h"
#include "freertos/queue.h"
#include "esp_log.h"
#include "config.h"
#include "http_client.h"
#include "cJSON.h"
#include "time.h"

extern QueueHandle_t event_queue;

static void add_timestamp(cJSON *obj)
{
  time_t now;
  struct tm timeinfo;
  char timestamp[64];

  time(&now);
  localtime_r(&now, &timeinfo);
  strftime(timestamp, sizeof(timestamp), "%Y-%m-%d %H:%M:%S", &timeinfo);
  cJSON_AddStringToObject(obj, "recorded_at", timestamp);
}

void telemetry_task(void *pvParameters)
{
  while (1)
  {
    vTaskDelay(pdMS_TO_TICKS(TELEMETRY_TASK_FREQUENCY_MS));

    cJSON *root = cJSON_CreateObject();
    cJSON *metrics = cJSON_AddArrayToObject(root, "metrics");

    event_t event;
    while (xQueueReceive(event_queue, &event, 0))
    {
      cJSON *car_light = cJSON_CreateObject();
      cJSON_AddStringToObject(car_light, "type", "car_light");
      cJSON_AddStringToObject(car_light, "unit", "quality");
      cJSON_AddStringToObject(car_light, "value", (event == EVENT_STATE_CHANGE) ? "GREEN" : "RED");
      add_timestamp(car_light);
      cJSON_AddItemToArray(metrics, car_light);

      cJSON *ped_light = cJSON_CreateObject();
      cJSON_AddStringToObject(ped_light, "type", "pedestrian_light");
      cJSON_AddStringToObject(ped_light, "unit", "quality");
      cJSON_AddStringToObject(ped_light, "value", (event == EVENT_STATE_CHANGE) ? "GREEN" : "RED");
      add_timestamp(ped_light);
      cJSON_AddItemToArray(metrics, ped_light);

      cJSON *ped_button = cJSON_CreateObject();
      cJSON_AddStringToObject(ped_button, "type", "pedestrian_button");
      cJSON_AddStringToObject(ped_button, "unit", "quality");
      cJSON_AddStringToObject(ped_button, "value", (event == EVENT_BUTTON_PRESS) ? "ACTIVE" : "INACTIVE");
      add_timestamp(ped_button);
      cJSON_AddItemToArray(metrics, ped_button);
    }

    char *json_str = cJSON_PrintUnformatted(root);
    http_post(TELEMETRY_URL, json_str);
    cJSON_Delete(root);
    free(json_str);
  }
}
