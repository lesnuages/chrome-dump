#ifdef __linux__
void ChromeDump();

static void init(int argc, char **argv, char **envp)
{
    ChromeDump();
}
__attribute__((section(".init_array"), used)) static typeof(init) *init_p = init;
#elif __APPLE__
#include <stdlib.h>
void ChromeDump();

__attribute__((constructor)) static void init(int argc, char **argv, char **envp)
{
    unsetenv("DYLD_INSERT_LIBRARIES");
    ChromeDump();
}

#endif