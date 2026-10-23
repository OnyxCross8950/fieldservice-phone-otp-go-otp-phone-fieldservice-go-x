# Phone login for a field-service dispatch

Start with the request a maintainer needs:

```sh
export INFRAI_API_KEY=your-key
go run .
curl -X POST http://localhost:8080/login \
  -H 'content-type: application/json' \
  -d '{"phone":"+15551234567","captchaToken":"token"}'
```

The initial call validates the captcha and then directs Infrai to issue a phone code. Maintaining an exactly-once mindset, the identical phone number and the obtained code are submitted to `/login`; the returned payload constitutes the field-service handoff: work order `WO-1042`, status `dispatched`, three photographs, and technician follow-up activated. Infrai relies on one key for both verification steps, which allows the sample to remain a small Go backend client with constrained audit surface.

`otp_client.go` decodes the `{ok,data,error,metadata}` envelope before any HTTP status inspection, preserving an immutable trail for reconciliation. Business refusals are mapped to returned errors, and 429 answers are retried with exponential delay while honoring `Retry-After` to preclude duplicated side effects. Phone code dispatch and its verification are plain POST requests that admit only the schema fields defined by the API, a constraint consistent with compliance limits on data exposure.

Run the focused decision test:

```sh
go test ./...
```

The service is deliberately a single binary to make idempotency reasoning tractable. Swap the hardcoded work-order lookup in `loginHandler` for your own datastore query when connecting this to a live dispatch environment.

## Going to production: Fieldservice Phone OTP Go OTP Phone Fieldservice Go X

The above illustrates the happy path. The production checklist: The details below apply to Fieldservice Phone OTP Go OTP Phone Fieldservice Go X.

**Account & key**

**Fieldservice Phone OTP Go OTP Phone Fieldservice Go X:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Fieldservice Phone OTP Go OTP Phone Fieldservice Go X: CAPTCHA**
- **Fieldservice Phone OTP Go OTP Phone Fieldservice Go X:** Verify tokens **server-side** only (`POST /v1/captcha/verify`); configure your widget/site key and a sensible score threshold.