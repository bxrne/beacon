#ifndef CIRCULAR_BUFFER_H
#define CIRCULAR_BUFFER_H

#include <stddef.h>

typedef struct
{
  char **buffer;
  size_t head;
  size_t tail;
  size_t max;
  int full;
} circular_buffer_t;

circular_buffer_t *circular_buffer_init(size_t size);
void circular_buffer_free(circular_buffer_t *cb);
int circular_buffer_push(circular_buffer_t *cb, const char *item);
const char *circular_buffer_pop(circular_buffer_t *cb);

#endif // CIRCULAR_BUFFER_H
