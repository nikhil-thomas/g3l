**note** on gke got error

```bash
panic: no Auth Provider found for name "gcp"
```

👉 fix: import `_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"`