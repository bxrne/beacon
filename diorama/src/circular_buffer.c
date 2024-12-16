#include "circular_buffer.h"
#include <stdlib.h>
#include <string.h>

circular_buffer_t *circular_buffer_init(size_t size)
{
  circular_buffer_t *cb = (circular_buffer_t *)malloc(sizeof(circular_buffer_t));
  cb->buffer = (char **)malloc(size * sizeof(char *));
  cb->max = size;
  cb->head = 0;
  cb->tail = 0;
  cb->full = 0;
  return cb;
}

void circular_buffer_free(circular_buffer_t *cb)
{
  for (size_t i = 0; i < cb->max; i++)
  {
    free(cb->buffer[i]);
  }
  free(cb->buffer);
  free(cb);
}

int circular_buffer_push(circular_buffer_t *cb, const char *item)
{
  size_t len = strlen(item) + 1;
  cb->buffer[cb->head] = (char *)malloc(len);
  strncpy(cb->buffer[cb->head], item, len);

  if (cb->full)
  {
    cb->tail = (cb->tail + 1) % cb->max;
  }

  cb->head = (cb->head + 1) % cb->max;
  cb->full = (cb->head == cb->tail);

  return 0;
}

const char *circular_buffer_pop(circular_buffer_t *cb)
{
  if (cb->head == cb->tail && !cb->full)
  {
    return NULL;
  }

  const char *item = cb->buffer[cb->tail];
  cb->tail = (cb->tail + 1) % cb->max;
  cb->full = 0;

  return item;
}
