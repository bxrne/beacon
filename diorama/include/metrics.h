#ifndef METRICS_H
#define METRICS_H

#include <stddef.h>
#include "freertos/FreeRTOS.h"
#include "freertos/semphr.h"

typedef enum
{
  LIGHT_GREEN,
  LIGHT_YELLOW,
  LIGHT_RED
} LightColor;

// Initialize the metrics buffers
void init_metrics_buffers(size_t size);

// Add light states to the circular buffers
void add_car_light_state(LightColor state);
void add_ped_light_state(LightColor state);

// Retrieve the most recent light states
LightColor get_recent_car_light_state();
LightColor get_recent_ped_light_state();

// Convert light color to string
const char *light_color_to_string(LightColor color);

// Free the metrics buffers
void free_metrics_buffers();

#endif // METRICS_H
