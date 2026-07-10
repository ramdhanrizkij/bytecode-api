# Storage

## Filesystem

Local storage is the default provider. It creates the root path and bucket directories at startup.

Default values:

- `STORAGE_LOCAL_PATH=storage`
- `STORAGE_BASE_URL=/storage`
- `STORAGE_DEFAULT_BUCKET=uploads`

When local storage is active, the server exposes static files at:

```text
<STORAGE_BASE_URL>/*
```

## S3

Amazon S3 is not implemented directly.

## MinIO

MinIO is supported through `github.com/minio/minio-go/v7` when:

```bash
STORAGE_PROVIDER=minio
```

The provider checks configured buckets and creates missing buckets.

## Local Storage

Local object writes use `os.Create` and `io.Copy`. Object paths are resolved under the configured root and checked to prevent path traversal outside that root.

## Uploads

Generic `Provider.Put` is implemented, but no HTTP upload endpoint is present in the analyzed codebase. User profile pictures are stored as object keys supplied in request bodies.

## Static Assets

Local storage files are served by Fiber static middleware only when the active provider is `local`.

## Object URL Format

| Provider | URL |
| --- | --- |
| local | `<base_url>/<bucket>/<escaped_key>` |
| minio | `<MINIO_PUBLIC_URL>/<bucket>/<escaped_key>` when public URL is set |
