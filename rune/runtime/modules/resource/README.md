# Resource module

The resource module provides provider-neutral CPU and memory validation and
host inspection. It does not reserve resources or modify virtual machines.

```bash
rune resource host
rune resource usage 1234
rune resource validate --cpu 4 --memory 4096
```

Providers such as Firecracker consume `resource::validate_cpu` and
`resource::validate_memory` through the public API.
