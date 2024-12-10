#include "time_sync.h"
#include "esp_sntp.h"
#include "esp_log.h"

static const char *TAG = "TIME_SYNC";

esp_err_t initialize_sntp(void)
{
  sntp_setoperatingmode(SNTP_OPMODE_POLL);
  sntp_setservername(0, "pool.ntp.org");
  sntp_init();

  // Wait for time to be set
  int retry = 0;
  const int retry_count = 10;
  while (sntp_get_sync_status() == SNTP_SYNC_STATUS_RESET && ++retry < retry_count)
  {
    ESP_LOGI(TAG, "Waiting for system time to be set... (%d/%d)", retry, retry_count);
    vTaskDelay(pdMS_TO_TICKS(2000));
  }

  if (retry == retry_count)
  {
    ESP_LOGE(TAG, "Failed to sync time");
    return ESP_FAIL;
  }

  return ESP_OK;
}
