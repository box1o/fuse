setup()
{
    PROJECT_ROOT="$(cd -- "${BATS_TEST_DIRNAME}/.." && pwd -P)"
    source "${PROJECT_ROOT}/lib/init.sh"
}
