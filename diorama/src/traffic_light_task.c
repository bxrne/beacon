#include <stdint.h>
#include <stdbool.h>
#include <stddef.h>
#include <time.h> // Add this header

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

static esp_http_client_handle_t client = NULL;
static SemaphoreHandle_t http_mutex = NULL;

static esp_err_t _http_event_handler(esp_http_client_event_t *evt)
{
  switch (evt->event_id)
  {
  case HTTP_EVENT_ERROR:
    ESP_LOGW("HTTP_CLIENT", "HTTP_EVENT_ERROR");
    break;
  case HTTP_EVENT_DISCONNECTED:
    ESP_LOGW("HTTP_CLIENT", "HTTP_EVENT_DISCONNECTED");
    break;
  default:
    break;
  }
  return ESP_OK;
}

static void init_http_client(void)
{
  if (client == NULL)
  {
    esp_http_client_config_t config = {
        .url = API_URL,
        .event_handler = _http_event_handler,
        .timeout_ms = 5000,
        .keep_alive_enable = true,
        .transport_type = HTTP_TRANSPORT_OVER_TCP, // Changed from SSL to TCP
        .buffer_size = 2048,
        .buffer_size_tx = 2048,
    };

    // Initialize the client
    client = esp_http_client_init(&config);
    if (client == NULL)
    {
      ESP_LOGE("HTTP_CLIENT", "Failed to initialize HTTP client");
      return;
    }

    // Create mutex if not exists
    if (http_mutex == NULL)
    {
      http_mutex = xSemaphoreCreateMutex();
    }
  }
}

static char *get_utc_timestamp(void)
{
  static char timestamp[25]; // Buffer for ISO 8601 timestamp
  time_t now;
  struct tm timeinfo;

  time(&now);
  gmtime_r(&now, &timeinfo);
  strftime(timestamp, sizeof(timestamp), "%Y-%m-%dT%H:%M:%SZ", &timeinfo);

  return timestamp;
}

static void send_metric_with_type(const char *type, const char *value, const char *unit)
{
  if (client == NULL)
  {
    init_http_client();
    if (client == NULL)
    {
      ESP_LOGE("HTTP_CLIENT", "Could not initialize client");
      return;
    }
  }

  if (xSemaphoreTake(http_mutex, pdMS_TO_TICKS(5000)) != pdTRUE)
  {
    ESP_LOGE("HTTP_CLIENT", "Could not acquire mutex");
    return;
  }

  // Create JSON payload
  cJSON *root = cJSON_CreateObject();
  cJSON *metrics_array = cJSON_AddArrayToObject(root, "metrics");
  cJSON *metric = cJSON_CreateObject();
  cJSON_AddItemToArray(metrics_array, metric);
  cJSON_AddStringToObject(metric, "type", type);
  cJSON_AddStringToObject(metric, "value", value);
  cJSON_AddStringToObject(metric, "unit", unit);
  cJSON_AddStringToObject(metric, "recorded_at", get_utc_timestamp());

  char *post_data = cJSON_PrintUnformatted(root);
  cJSON_Delete(root);

  // Set headers
  esp_http_client_set_header(client, "Content-Type", "application/json");
  esp_http_client_set_header(client, "X-DeviceID", DEVICE_ID);
  esp_http_client_set_header(client, "Connection", "keep-alive");

  // Set post data
  esp_http_client_set_method(client, HTTP_METHOD_POST);
  esp_http_client_set_post_field(client, post_data, strlen(post_data));

  // Perform request with connection check
  esp_err_t err = esp_http_client_perform(client);
  if (err == ESP_OK)
  {
    int status_code = esp_http_client_get_status_code(client);
    ESP_LOGI("HTTP_CLIENT", "HTTP POST Status = %d", status_code);
  }
  else
  {
    ESP_LOGE("HTTP_CLIENT", "HTTP POST request failed: %s", esp_err_to_name(err));
    // Clean up and recreate client on error
    esp_http_client_cleanup(client);
    client = NULL;
    init_http_client(); // Try to reinitialize for next time
  }

  free(post_data);
  xSemaphoreGive(http_mutex);
}

static void send_traffic_light_metric(const char *color)
{
  send_metric_with_type("traffic_light", color, "color");
}

static void send_crossing_light_metric(const char *color)
{
  send_metric_with_type("crossing_light", color, "color");
}

static void send_crossing_button_metric(bool pressed)
{
  send_metric_with_type("crossing_button", pressed ? "true" : "false", "bool");
}

void traffic_light_task(void *pvParameters)
{
  // Initialize HTTP client first
  init_http_client();

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
        send_crossing_button_metric(true);
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
      send_traffic_light_metric("green");
      send_crossing_light_metric("red");
      ESP_LOGI("STATE", "Car light: Green");

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
      send_traffic_light_metric("yellow");
      ESP_LOGI("STATE", "Car light: Yellow");
      vTaskDelay(pdMS_TO_TICKS(CAR_YELLOW_DURATION));
      current_state = STATE_CAR_RED;
      break;

    case STATE_CAR_RED:
      gpio_set_level(CAR_YELLOW_PIN, 0);
      gpio_set_level(CAR_RED_PIN, 1);
      send_traffic_light_metric("red");
      ESP_LOGI("STATE", "Car light: Red");

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
      send_crossing_light_metric("green");
      ESP_LOGI("STATE", "Pedestrian light: Green");
      vTaskDelay(pdMS_TO_TICKS(PED_GREEN_DURATION));

      gpio_set_level(PED_GREEN_PIN, 0);
      gpio_set_level(PED_RED_PIN, 1);
      send_crossing_light_metric("red");
      send_crossing_button_metric(false);
      pedestrian_button = false;
      current_state = STATE_CAR_GREEN;
      break;
    }
  }
}