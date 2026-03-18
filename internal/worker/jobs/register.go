package jobs

import (
	identityJob "github.com/ramdhanrizki/bytecode-api/internal/identity/application/job"
	sharedQueue "github.com/ramdhanrizki/bytecode-api/internal/shared/queue"
	workerQueue "github.com/ramdhanrizki/bytecode-api/internal/worker/queue"
)

type Dependencies struct {
	EmailVerificationHandler sharedQueue.Handler
}

func Register(registry *workerQueue.Registry, deps Dependencies) {
	if deps.EmailVerificationHandler != nil {
		registry.Register(identityJob.EmailVerificationJobName, deps.EmailVerificationHandler)
	}
}
