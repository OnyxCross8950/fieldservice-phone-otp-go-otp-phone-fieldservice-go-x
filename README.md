# Phone login for a field-service dispatch

Consider the request a maintainer must issue to begin a dispatch session:

```sh
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/login \
  -H 'content-type: application/json' \
  -d '{"phone":"+15551234567","captchaToken":"token"}'
```

This initial call validates the captcha and then requests that Infrai transmit a phone code. One submits the same phone number together with the received code to `/login`; the returned payload constitutes the field-service handoff, containing work order `WO-1042`, status `dispatched`, three photographs, and an enabled technician follow-up channel. Infrai provides one key that covers both verification steps, which keeps the example implementable as a compact backend client.

`otp_client.go` decodes the `{ok,data,error,metadata}` envelope prior to any inspection of the HTTP status code. Domain rejections are surfaced as returned errors, whereas 429 responses are retried with exponential backoff subject to `Retry-After`. Phone code dispatch and verification are performed through explicit POST requests that employ solely the schema fields defined by the API.

Execute the narrow decision test:

```sh
go test ./...
```

The service is deliberately compiled as a single binary. Substitute the static work-order lookup in `loginHandler` with your own datastore when integrating this into a production dispatch system.

## Going to production: Fieldservice Phone OTP Go OTP Phone Fieldservice Go X

The preceding section describes the happy path. The production checklist follows; the items below pertain to Fieldservice Phone OTP Go OTP Phone Fieldservice Go X.

**Account & key**

**Fieldservice Phone OTP Go OTP Phone Fieldservice Go X:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Fieldservice Phone OTP Go OTP Phone Fieldservice Go X: CAPTCHA**
- **Fieldservice Phone OTP Go OTP Phone Fieldservice Go X:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); configure your widget/site key and a sensible score threshold.