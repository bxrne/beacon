#ifndef CONFIG_H
#define CONFIG_H

#define TAG "TRAFFIC_LIGHT"

#define CAR_GREEN_PIN 18
#define CAR_GREEN_DURATION 5000 // 5 seconds
#define CAR_YELLOW_PIN 19
#define CAR_YELLOW_DURATION 1000 // 1 second
#define CAR_RED_PIN 21
#define CAR_RED_DURATION 3000 // 3 seconds

#define PED_GREEN_PIN 22
#define PED_GREEN_DURATION 5000 // 5 seconds
#define PED_RED_PIN 23
#define PED_BUTTON_PIN 25

#define BUTTON_PIN_CHECK_INTERVAL 50 // 50 ms
#define DEBOUNCE_THRESHOLD 3         // 150 ms (3 * 50ms)
#define BUTTON_PRESS_THRESHOLD 20    // 1 second (20 * 50ms)

#define DEBOUNCE_TIME_MS 50 // Debounce time in milliseconds
#define GPIO_INTERRUPT_TRIGGER GPIO_INTR_NEGEDGE
#define BUTTON_QUEUE_SIZE 5 // Size of button event queue

typedef enum
{
  EVENT_BUTTON_PRESS
} event_t;

typedef enum
{
  BUTTON_RELEASED,
  BUTTON_PRESSED,
  BUTTON_DEBOUNCING
} button_state_t;

#endif // CONFIG_H
