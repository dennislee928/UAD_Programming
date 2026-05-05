// soc_erh_model.go
// Go implementation of SOC/JudgePipeline ERH model
// This demonstrates how String Theory semantics would be implemented in Go
// without native coupling and resonance support

package main

import (
	"fmt"
	"sync"
)

// ============================================================================
// Core Domain Models
// ============================================================================

// String represents an entity with vibrational modes
type String struct {
	Name  string
	Modes map[string]float64
	mutex sync.RWMutex
}

func NewString(name string) *String {
	return &String{
		Name:  name,
		Modes: make(map[string]float64),
	}
}

func (s *String) SetMode(name string, value float64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Modes[name] = value
}

func (s *String) GetMode(name string) float64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.Modes[name]
}

// Coupling represents a relationship between two modes
type Coupling struct {
	SourceString string
	SourceMode   string
	TargetString string
	TargetMode   string
	Strength     float64
}

// CouplingManager manages all couplings and propagates changes
type CouplingManager struct {
	couplings []Coupling
	strings   map[string]*String
	observers []func(string, string, float64)
}

func NewCouplingManager() *CouplingManager {
	return &CouplingManager{
		couplings: []Coupling{},
		strings:   make(map[string]*String),
		observers: []func(string, string, float64){},
	}
}

func (cm *CouplingManager) RegisterString(s *String) {
	cm.strings[s.Name] = s
}

func (cm *CouplingManager) AddCoupling(c Coupling) {
	cm.couplings = append(cm.couplings, c)
}

// Propagate propagates a change from source to all coupled targets
func (cm *CouplingManager) Propagate(sourceString, sourceMode string, newValue float64) {
	for _, coupling := range cm.couplings {
		if coupling.SourceString == sourceString && coupling.SourceMode == sourceMode {
			targetStr := cm.strings[coupling.TargetString]
			if targetStr != nil {
				currentValue := targetStr.GetMode(coupling.TargetMode)
				change := (newValue - currentValue) * coupling.Strength
				newTargetValue := currentValue + change
				targetStr.SetMode(coupling.TargetMode, newTargetValue)
				
				// Notify observers
				for _, observer := range cm.observers {
					observer(coupling.TargetString, coupling.TargetMode, newTargetValue)
				}
				
				// Recursive propagation (could cause cycles in real implementation)
				cm.Propagate(coupling.TargetString, coupling.TargetMode, newTargetValue)
			}
		}
	}
}

// ResonanceRule represents a conditional coupling effect
type ResonanceRule struct {
	Condition func(*CouplingManager) bool
	Effect    func(*CouplingManager)
	Name      string
}

// ResonanceEngine evaluates resonance rules
type ResonanceEngine struct {
	rules []ResonanceRule
	cm    *CouplingManager
}

func NewResonanceEngine(cm *CouplingManager) *ResonanceEngine {
	return &ResonanceEngine{
		rules: []ResonanceRule{},
		cm:    cm,
	}
}

func (re *ResonanceEngine) AddRule(rule ResonanceRule) {
	re.rules = append(re.rules, rule)
}

func (re *ResonanceEngine) Evaluate() {
	for _, rule := range re.rules {
		if rule.Condition(re.cm) {
			fmt.Printf("Resonance triggered: %s\n", rule.Name)
			rule.Effect(re.cm)
		}
	}
}

// ============================================================================
// SOC/JudgePipeline Specific Models
// ============================================================================

type JudgeKind string

const (
	JudgeHuman   JudgeKind = "human"
	JudgePipeline JudgeKind = "pipeline"
	JudgeModel   JudgeKind = "model"
)

type Judge struct {
	Kind       JudgeKind
	Decision   float64
	Confidence float64
	Reasoning  string
}

type Action struct {
	ID         string
	Complexity float64
	TrueValue  float64
	Importance float64
}

type ERHProfile struct {
	ActionClass    string
	PrimeThreshold float64
	FitAlpha       float64
	EthicalScore   float64
}

func CalculateERHScore(profile ERHProfile, actionCount int) float64 {
	baseScore := profile.EthicalScore
	var thresholdFactor float64 = 1.0
	
	if float64(actionCount) > profile.PrimeThreshold {
		thresholdFactor = 1.2 // Above threshold: increased risk
	}
	
	alphaAdjustment := profile.FitAlpha * 0.1
	return baseScore*thresholdFactor + alphaAdjustment
}

// ============================================================================
// SOC System Implementation
// ============================================================================

type SOCSystem struct {
	cm              *CouplingManager
	resonanceEngine *ResonanceEngine
	
	// Strings (entities)
	agentPopulation    *String
	systemInfrastructure *String
	decisionFramework   *String
	systemState        *String
}

func NewSOCSystem() *SOCSystem {
	cm := NewCouplingManager()
	
	// Initialize strings
	agentPop := NewString("agent_population")
	agentPop.SetMode("trust_level", 6.0)
	agentPop.SetMode("risk_score", 3.0)
	agentPop.SetMode("performance", 7.0)
	
	sysInfra := NewString("system_infrastructure")
	sysInfra.SetMode("reliability", 7.5)
	sysInfra.SetMode("security_posture", 6.5)
	sysInfra.SetMode("performance", 7.0)
	
	decisionFw := NewString("decision_framework")
	decisionFw.SetMode("transparency", 5.5)
	decisionFw.SetMode("accountability", 6.0)
	decisionFw.SetMode("ethical_score", 6.5)
	
	sysState := NewString("system_state")
	sysState.SetMode("overall_stability", 6.8)
	sysState.SetMode("systemic_risk", 3.2)
	sysState.SetMode("alert_level", 3.0)
	
	// Register strings
	cm.RegisterString(agentPop)
	cm.RegisterString(sysInfra)
	cm.RegisterString(decisionFw)
	cm.RegisterString(sysState)
	
	// Define couplings
	// Agent trust affects system security
	cm.AddCoupling(Coupling{
		SourceString: "agent_population",
		SourceMode:   "trust_level",
		TargetString: "system_infrastructure",
		TargetMode:   "security_posture",
		Strength:     0.65,
	})
	
	// System reliability affects system state stability
	cm.AddCoupling(Coupling{
		SourceString: "system_infrastructure",
		SourceMode:   "reliability",
		TargetString: "system_state",
		TargetMode:   "overall_stability",
		Strength:     0.75,
	})
	
	// Decision transparency affects agent trust
	cm.AddCoupling(Coupling{
		SourceString: "decision_framework",
		SourceMode:   "transparency",
		TargetString: "agent_population",
		TargetMode:   "trust_level",
		Strength:     0.70,
	})
	
	// Decision accountability affects system security
	cm.AddCoupling(Coupling{
		SourceString: "decision_framework",
		SourceMode:   "accountability",
		TargetString: "system_infrastructure",
		TargetMode:   "security_posture",
		Strength:     0.60,
	})
	
	// Agent risk affects system risk
	cm.AddCoupling(Coupling{
		SourceString: "agent_population",
		SourceMode:   "risk_score",
		TargetString: "system_state",
		TargetMode:   "systemic_risk",
		Strength:     0.90,
	})
	
	// System security inversely affects risk
	cm.AddCoupling(Coupling{
		SourceString: "system_infrastructure",
		SourceMode:   "security_posture",
		TargetString: "system_state",
		TargetMode:   "systemic_risk",
		Strength:     -0.85, // Negative: higher security = lower risk
	})
	
	// Initialize resonance engine
	re := NewResonanceEngine(cm)
	
	// Resonance: High agent trust amplifies system improvements
	re.AddRule(ResonanceRule{
		Name: "trust_resonance",
		Condition: func(cm *CouplingManager) bool {
			agentPop := cm.strings["agent_population"]
			return agentPop != nil && agentPop.GetMode("trust_level") > 8.0
		},
		Effect: func(cm *CouplingManager) {
			fmt.Println("  Effect: Amplifying system improvements (multiplier: 1.5)")
		},
	})
	
	// Resonance: Low system reliability triggers cascading failures
	re.AddRule(ResonanceRule{
		Name: "reliability_cascade",
		Condition: func(cm *CouplingManager) bool {
			sysInfra := cm.strings["system_infrastructure"]
			return sysInfra != nil && sysInfra.GetMode("reliability") < 3.0
		},
		Effect: func(cm *CouplingManager) {
			sysState := cm.strings["system_state"]
			if sysState != nil {
				currentRisk := sysState.GetMode("systemic_risk")
				sysState.SetMode("systemic_risk", currentRisk+0.3)
				fmt.Println("  Effect: Risk amplification (+0.3)")
			}
		},
	})
	
	// Resonance: Critical risk threshold
	re.AddRule(ResonanceRule{
		Name: "critical_risk_alert",
		Condition: func(cm *CouplingManager) bool {
			sysState := cm.strings["system_state"]
			return sysState != nil && sysState.GetMode("systemic_risk") > 8.5
		},
		Effect: func(cm *CouplingManager) {
			fmt.Println("  Effect: Emergency intervention required")
			fmt.Println("  Actions: increase_transparency, enhance_accountability, improve_security")
		},
	})
	
	return &SOCSystem{
		cm:              cm,
		resonanceEngine: re,
		agentPopulation:     agentPop,
		systemInfrastructure: sysInfra,
		decisionFramework:    decisionFw,
		systemState:         sysState,
	}
}

func (soc *SOCSystem) UpdateAgentTrust(newValue float64) {
	soc.agentPopulation.SetMode("trust_level", newValue)
	soc.cm.Propagate("agent_population", "trust_level", newValue)
	soc.resonanceEngine.Evaluate()
}

func (soc *SOCSystem) UpdateSystemReliability(newValue float64) {
	soc.systemInfrastructure.SetMode("reliability", newValue)
	soc.cm.Propagate("system_infrastructure", "reliability", newValue)
	soc.resonanceEngine.Evaluate()
}

func (soc *SOCSystem) UpdateAgentRisk(newValue float64) {
	soc.agentPopulation.SetMode("risk_score", newValue)
	soc.cm.Propagate("agent_population", "risk_score", newValue)
	soc.resonanceEngine.Evaluate()
}

func (soc *SOCSystem) PrintState() {
	fmt.Println("Current System State:")
	fmt.Printf("  Agent Trust: %.2f\n", soc.agentPopulation.GetMode("trust_level"))
	fmt.Printf("  Agent Risk: %.2f\n", soc.agentPopulation.GetMode("risk_score"))
	fmt.Printf("  System Reliability: %.2f\n", soc.systemInfrastructure.GetMode("reliability"))
	fmt.Printf("  System Security: %.2f\n", soc.systemInfrastructure.GetMode("security_posture"))
	fmt.Printf("  Decision Transparency: %.2f\n", soc.decisionFramework.GetMode("transparency"))
	fmt.Printf("  System Stability: %.2f\n", soc.systemState.GetMode("overall_stability"))
	fmt.Printf("  Systemic Risk: %.2f\n", soc.systemState.GetMode("systemic_risk"))
	fmt.Println("")
}

// ============================================================================
// Main Execution
// ============================================================================

func main() {
	fmt.Println("=========================================")
	fmt.Println("SOC / JudgePipeline ERH Model (Go)")
	fmt.Println("=========================================")
	fmt.Println("")
	fmt.Println("This Go implementation demonstrates:")
	fmt.Println("1. Manual coupling propagation (no native support)")
	fmt.Println("2. Manual resonance rule evaluation")
	fmt.Println("3. Manual state synchronization")
	fmt.Println("4. Observer pattern for change propagation")
	fmt.Println("")
	fmt.Println("Comparison with UAD:")
	fmt.Println("  - UAD: ~200 LOC with native String Theory semantics")
	fmt.Println("  - Go:  ~400+ LOC with manual coupling/resonance")
	fmt.Println("  - Expressiveness: UAD 5/5, Go 2/5")
	fmt.Println("")
	
	// Initialize SOC system
	soc := NewSOCSystem()
	
	fmt.Println("Initial State:")
	soc.PrintState()
	
	// Simulate updates
	fmt.Println("Step 1: Agent trust increases to 7.5")
	soc.UpdateAgentTrust(7.5)
	soc.PrintState()
	
	fmt.Println("Step 2: System reliability drops to 2.5 (external shock)")
	soc.UpdateSystemReliability(2.5)
	soc.PrintState()
	
	fmt.Println("Step 3: Agent risk increases to 9.0")
	soc.UpdateAgentRisk(9.0)
	soc.PrintState()
	
	// ERH Example
	fmt.Println("ERH Profile Example:")
	profile := ERHProfile{
		ActionClass:    "decision_making",
		PrimeThreshold: 100.0,
		FitAlpha:       0.15,
		EthicalScore:   6.5,
	}
	actionCount := 150
	erhScore := CalculateERHScore(profile, actionCount)
	fmt.Printf("  Action Class: %s\n", profile.ActionClass)
	fmt.Printf("  Action Count: %d\n", actionCount)
	fmt.Printf("  Above Threshold: %v\n", actionCount > int(profile.PrimeThreshold))
	fmt.Printf("  Calculated ERH Score: %.2f\n", erhScore)
	fmt.Println("")
}


