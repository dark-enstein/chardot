package agent

import (
	"context"
	"fmt"
)

func GetSpeedFromCtx(ctx context.Context, speedType MovType) (*Speed, error) {
	speed, ok := ctx.Value(speedType).(*Speed)
	if !ok {
		return nil, fmt.Errorf("speed of type %v not found in context", speedType)
	}
	return speed, nil
}

func GetAgentFromCtx(ctx context.Context) (Agent, error) {
	agent, ok := ctx.Value(AGENT).(Agent)
	if !ok {
		return nil, fmt.Errorf("agent not found in context")
	}
	return agent, nil
}
