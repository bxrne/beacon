#ifndef CONFIG_H
#define CONFIG_H

#include <stdint.h>
#include <stdbool.h>

// Traffic light GPIO pins and durations
#define CAR_GREEN_PIN 18
#define CAR_GREEN_DURATION 5000 // Milliseconds
#define CAR_YELLOW_PIN 19
#define CAR_YELLOW_DURATION 1000
#define CAR_RED_PIN 21
#define CAR_RED_DURATION 3000

#define PED_GREEN_PIN 22
#define PED_GREEN_DURATION 5000
#define PED_RED_PIN 23
#define PED_BUTTON_PIN 25

// Button configurations
#define DEBOUNCE_TIME_MS 50
#define EVENT_QUEUE_SIZE 10
#define BUTTON_TASK_STACK_SIZE 2048
#define BUTTON_TASK_PRIORITY 10
#define TRAFFIC_LIGHT_TASK_STACK_SIZE 4096
#define TRAFFIC_LIGHT_TASK_PRIORITY 5

#define ESP_INTR_FLAG_DEFAULT 0

// Circular buffer size
#define CIRCULAR_BUFFER_SIZE 1024

// Define event types
typedef enum
{
  EVENT_BUTTON_PRESS,
  EVENT_REBOOT_COMMAND
} event_t;

#endif // CONFIG_H
