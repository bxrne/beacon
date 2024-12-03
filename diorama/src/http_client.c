#include "http_client.h"
#include "config.h"
#include "esp_log.h"
#include "esp_http_client.h"
#include "esp_system.h"
#include "esp_efuse_table.h"
#include "esp_efuse.h"

esp_err_t _http_event_handler(esp_http_client_event_t *evt)
{
  switch (evt->event_id)
  {
  case HTTP_EVENT_ERROR:
    ESP_LOGI("HTTP_EVENT", "HTTP_EVENT_ERROR");
    break;
  case HTTP_EVENT_ON_CONNECTED:
    ESP_LOGI("HTTP_EVENT", "HTTP_EVENT_ON_CONNECTED");
    break;
  case HTTP_EVENT_HEADER_SENT:
    ESP_LOGI("HTTP_EVENT", "HTTP_EVENT_HEADER_SENT");
    break;
  case HTTP_EVENT_ON_HEADER:
    ESP_LOGI("HTTP_EVENT", "HTTP_EVENT_ON_HEADER, key=%s, value=%s", evt->header_key, evt->header_value);
    break;
  case HTTP_EVENT_ON_DATA:
    ESP_LOGI("HTTP_EVENT", "HTTP_EVENT_ON_DATA, len=%d", evt->data_len);
    break;
  case HTTP_EVENT_ON_FINISH:
    ESP_LOGI("HTTP_EVENT", "HTTP_EVENT_ON_FINISH");
    break;
  case HTTP_EVENT_DISCONNECTED:
    ESP_LOGI("HTTP_EVENT", "HTTP_EVENT_DISCONNECTED");
    break;
  case HTTP_EVENT_REDIRECT:
    ESP_LOGI("HTTP_EVENT", "HTTP_EVENT_REDIRECT");
    break;
  }
  return ESP_OK;
}

esp_err_t http_post(const char *url, const char *data)
{
  esp_http_client_config_t config = {
      .url = url,
      .event_handler = _http_event_handler,
  };
  esp_http_client_handle_t client = esp_http_client_init(&config);

  if (client == NULL)
  {
    ESP_LOGE("HTTP_POST", "Failed to initialize HTTP client");
    return ESP_FAIL;
  }

  esp_http_client_set_method(client, HTTP_METHOD_POST);
  esp_http_client_set_header(client, "Content-Type", "application/json");
  esp_http_client_set_header(client, "X-DeviceID", "esp");
  esp_http_client_set_post_field(client, data, strlen(data));

  esp_err_t err = esp_http_client_perform(client);
  if (err == ESP_OK)
  {
    ESP_LOGI("HTTP_POST", "HTTP POST Status = %d, content_length = %d",
             (int)esp_http_client_get_status_code(client),
             (int)esp_http_client_get_content_length(client));
  }
  else
  {
    ESP_LOGE("HTTP_POST", "HTTP POST request failed: %s", esp_err_to_name(err));
  }

  esp_http_client_cleanup(client);
  return err;
}