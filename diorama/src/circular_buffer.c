#include "circular_buffer.h"
#include <stdlib.h>
#include <string.h>

circular_buffer_t *circular_buffer_init(size_t size)
{
  circular_buffer_t *cb = (circular_buffer_t *)malloc(sizeof(circular_buffer_t));
  cb->buffer = (int *)malloc(size * sizeof(int));
  cb->max = size;
  cb->head = 0;
  cb->tail = 0;
  cb->full = 0;
  return cb;
}

void circular_buffer_free(circular_buffer_t *cb)
{
  free(cb->buffer);
  free(cb);
}

int circular_buffer_push(circular_buffer_t *cb, int item)
{
  cb->buffer[cb->head] = item;
  cb->head = (cb->head + 1) % cb->max;

  if (cb->full)
  {
    cb->tail = (cb->tail + 1) % cb->max;
  }

  cb->full = (cb->head == cb->tail);
  return 0;
}

int circular_buffer_pop(circular_buffer_t *cb)
{
  if (circular_buffer_is_empty(cb))
  {
    return -1;
  }

  int item = cb->buffer[cb->tail];
  cb->full = 0;
  cb->tail = (cb->tail + 1) % cb->max;
  return item;
}

bool circular_buffer_is_empty(circular_buffer_t *cb)
{
  return (!cb->full && (cb->head == cb->tail));
}

int circular_buffer_peek_last(circular_buffer_t *cb)
{
  if (circular_buffer_is_empty(cb))
  {
    return -1;
  }

  size_t index = (cb->head == 0) ? cb->max - 1 : cb->head - 1;
  return cb->buffer[index];
}
