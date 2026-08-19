package keyparams

import "raiseexception.dev/odin/src/shared/domain/odinerrors"

type KeyParams struct {
	algorithm   string
	salt        string
	iterations  int
	memory      int
	parallelism int
}

func New(algorithm string, iterations, memory, parallelism int, salt string) (KeyParams, error) {
	if algorithm == "" {
		return KeyParams{}, buildError("algorithm is required", "El algoritmo es obligatorio")
	}
	if iterations <= 0 {
		return KeyParams{}, buildError("iterations must be positive", "Las iteraciones deben ser positivas")
	}
	if memory <= 0 {
		return KeyParams{}, buildError("memory must be positive", "La memoria debe ser positiva")
	}
	if parallelism <= 0 {
		return KeyParams{}, buildError("parallelism must be positive", "El paralelismo debe ser positivo")
	}
	if salt == "" {
		return KeyParams{}, buildError("salt is required", "La sal es obligatoria")
	}
	return KeyParams{
		algorithm:   algorithm,
		salt:        salt,
		iterations:  iterations,
		memory:      memory,
		parallelism: parallelism,
	}, nil
}

func (self KeyParams) Algorithm() string {
	return self.algorithm
}

func (self KeyParams) Iterations() int {
	return self.iterations
}

func (self KeyParams) Memory() int {
	return self.memory
}

func (self KeyParams) Parallelism() int {
	return self.parallelism
}

func (self KeyParams) Salt() string {
	return self.salt
}

func buildError(message, external string) error {
	return odinerrors.NewErrorBuilder(message).
		WithExternalMessage(external).
		WithTag(odinerrors.Domain).
		Build()
}
