#ifndef CIRCULAR_BUFFER_H
#define CIRCULAR_BUFFER_H

#include <stddef.h>
#include <stdbool.h>

typedef struct
{
  int *buffer;
  size_t head;
  size_t tail;
  size_t max;
  bool full;
} circular_buffer_t;

circular_buffer_t *circular_buffer_init(size_t size);
void circular_buffer_free(circular_buffer_t *cb);
int circular_buffer_push(circular_buffer_t *cb, int item);
int circular_buffer_pop(circular_buffer_t *cb);

bool circular_buffer_is_empty(circular_buffer_t *cb);
int circular_buffer_peek_last(circular_buffer_t *cb);

#endif // CIRCULAR_BUFFER_H
