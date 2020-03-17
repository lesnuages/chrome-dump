#ifdef __linux__
void ChromeDump();

static void init(int argc, char **argv, char **envp)
{
    ChromeDump();
}
__attribute__((section(".init_array"), used)) static typeof(init) *init_p = init;
#endif