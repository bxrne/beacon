#include "metrics.h"
#include "circular_buffer.h"
#include "freertos/FreeRTOS.h"
#include "freertos/semphr.h"
#include "esp_log.h"
#include "globals.h"
#include "config.h"
#include "driver/gpio.h" // Include for gpio_get_level
#include <time.h>
#include <sys/time.h>

extern SemaphoreHandle_t xPedestrianSemaphore;

static circular_buffer_t *car_light_buffer;
static circular_buffer_t *ped_light_buffer;

void init_metrics_buffers(size_t size)
{
  car_light_buffer = circular_buffer_init(size);
  ped_light_buffer = circular_buffer_init(size);
}

void free_metrics_buffers()
{
  circular_buffer_free(car_light_buffer);
  circular_buffer_free(ped_light_buffer);
}

void update_light_buffers()
{
  const char *car_light_color;
  if (gpio_get_level(CAR_GREEN_PIN))
  {
    car_light_color = "GREEN";
  }
  else if (gpio_get_level(CAR_YELLOW_PIN))
  {
    car_light_color = "YELLOW";
  }
  else if (gpio_get_level(CAR_RED_PIN))
  {
    car_light_color = "RED";
  }
  else
  {
    car_light_color = "UNKNOWN";
  }
  circular_buffer_push(car_light_buffer, car_light_color);

  const char *ped_light_color;
  if (gpio_get_level(PED_GREEN_PIN))
  {
    ped_light_color = "GREEN";
  }
  else if (gpio_get_level(PED_RED_PIN))
  {
    ped_light_color = "RED";
  }
  else
  {
    ped_light_color = "UNKNOWN";
  }
  circular_buffer_push(ped_light_buffer, ped_light_color);
}

void get_metrics(char *buffer, size_t buffer_size)
{
  // Get the current time
  time_t now;
  time(&now);
  struct tm timeinfo;
  gmtime_r(&now, &timeinfo); // Get UTC time

  // Get the latest light colors from the buffers
  const char *car_light_color = circular_buffer_pop(car_light_buffer);
  const char *ped_light_color = circular_buffer_pop(ped_light_buffer);

  // Check if the light colors are NULL and set to the current GPIO level if they are
  if (!car_light_color)
  {
    if (gpio_get_level(CAR_GREEN_PIN))
    {
      car_light_color = "GREEN";
    }
    else if (gpio_get_level(CAR_YELLOW_PIN))
    {
      car_light_color = "YELLOW";
    }
    else if (gpio_get_level(CAR_RED_PIN))
    {
      car_light_color = "RED";
    }
    else
    {
      car_light_color = "UNKNOWN";
    }
  }

  if (!ped_light_color)
  {
    if (gpio_get_level(PED_GREEN_PIN))
    {
      ped_light_color = "GREEN";
    }
    else if (gpio_get_level(PED_RED_PIN))
    {
      ped_light_color = "RED";
    }
    else
    {
      ped_light_color = "UNKNOWN";
    }
  }

  // Format the metrics into the buffer
  snprintf(buffer, buffer_size, "Car Light: %s, Ped Light: %s, Time: %02d:%02d:%02d",
           car_light_color ? car_light_color : "UNKNOWN",
           ped_light_color ? ped_light_color : "UNKNOWN",
           timeinfo.tm_hour, timeinfo.tm_min, timeinfo.tm_sec);
}
