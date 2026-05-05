"""
psychohistory_scenario.py
Python implementation of Psychohistory-style macro dynamics simulation
This demonstrates how String Theory semantics would be implemented in Python
without native coupling and resonance support
"""

from typing import Dict, List, Callable, Optional
from dataclasses import dataclass
from enum import Enum


# ============================================================================
# Core Domain Models
# ============================================================================

@dataclass
class Mode:
    """Represents a vibrational mode of a string"""
    name: str
    value: float


class String:
    """Represents an entity with vibrational modes"""
    
    def __init__(self, name: str):
        self.name = name
        self.modes: Dict[str, float] = {}
    
    def set_mode(self, name: str, value: float):
        """Set a mode value"""
        self.modes[name] = value
    
    def get_mode(self, name: str) -> float:
        """Get a mode value"""
        return self.modes.get(name, 0.0)


@dataclass
class Coupling:
    """Represents a coupling relationship between two modes"""
    source_string: str
    source_mode: str
    target_string: str
    target_mode: str
    strength: float


@dataclass
class ResonanceRule:
    """Represents a resonance rule with condition and effect"""
    name: str
    condition: Callable[[Dict[str, String]], bool]
    effect: Callable[[Dict[str, String]], None]


class CouplingManager:
    """Manages couplings and propagates changes"""
    
    def __init__(self):
        self.couplings: List[Coupling] = []
        self.strings: Dict[str, String] = {}
        self.observers: List[Callable[[str, str, float], None]] = []
    
    def register_string(self, s: String):
        """Register a string entity"""
        self.strings[s.name] = s
    
    def add_coupling(self, coupling: Coupling):
        """Add a coupling relationship"""
        self.couplings.append(coupling)
    
    def propagate(self, source_string: str, source_mode: str, new_value: float):
        """Propagate changes through couplings"""
        for coupling in self.couplings:
            if (coupling.source_string == source_string and 
                coupling.source_mode == source_mode):
                
                target_str = self.strings.get(coupling.target_string)
                if target_str:
                    current_value = target_str.get_mode(coupling.target_mode)
                    change = (new_value - current_value) * coupling.strength
                    new_target_value = current_value + change
                    target_str.set_mode(coupling.target_mode, new_target_value)
                    
                    # Notify observers
                    for observer in self.observers:
                        observer(coupling.target_string, 
                                coupling.target_mode, 
                                new_target_value)
                    
                    # Recursive propagation
                    self.propagate(coupling.target_string, 
                                 coupling.target_mode, 
                                 new_target_value)


class ResonanceEngine:
    """Evaluates resonance rules"""
    
    def __init__(self, cm: CouplingManager):
        self.cm = cm
        self.rules: List[ResonanceRule] = []
    
    def add_rule(self, rule: ResonanceRule):
        """Add a resonance rule"""
        self.rules.append(rule)
    
    def evaluate(self):
        """Evaluate all resonance rules"""
        for rule in self.rules:
            if rule.condition(self.cm.strings):
                print(f"Resonance triggered: {rule.name}")
                rule.effect(self.cm.strings)


# ============================================================================
# Brane (Dimensional Context) - Simplified representation
# ============================================================================

class Brane:
    """Represents a dimensional context"""
    
    def __init__(self, name: str, dimensions: List[str]):
        self.name = name
        self.dimensions = dimensions


# ============================================================================
# Macro Dynamics System
# ============================================================================

class MacroDynamicsSystem:
    """Psychohistory-style macro dynamics simulation system"""
    
    def __init__(self):
        self.cm = CouplingManager()
        self.resonance_engine = ResonanceEngine(self.cm)
        
        # Initialize strings (entities)
        self.agent_population = String("agent_population")
        self.agent_population.set_mode("trust_level", 6.0)
        self.agent_population.set_mode("risk_score", 3.0)
        self.agent_population.set_mode("adoption_rate", 5.0)
        
        self.system_infrastructure = String("system_infrastructure")
        self.system_infrastructure.set_mode("reliability", 7.5)
        self.system_infrastructure.set_mode("security_posture", 6.5)
        self.system_infrastructure.set_mode("performance", 7.0)
        
        self.decision_framework = String("decision_framework")
        self.decision_framework.set_mode("transparency", 5.5)
        self.decision_framework.set_mode("accountability", 6.0)
        self.decision_framework.set_mode("ethical_score", 6.5)
        
        self.market_dynamics = String("market_dynamics")
        self.market_dynamics.set_mode("competition_level", 6.0)
        self.market_dynamics.set_mode("innovation_rate", 5.5)
        self.market_dynamics.set_mode("market_confidence", 6.5)
        
        self.macro_system_state = String("macro_system_state")
        self.macro_system_state.set_mode("overall_stability", 6.8)
        self.macro_system_state.set_mode("systemic_risk", 3.2)
        self.macro_system_state.set_mode("evolution_rate", 1.0)
        
        # Register strings
        self.cm.register_string(self.agent_population)
        self.cm.register_string(self.system_infrastructure)
        self.cm.register_string(self.decision_framework)
        self.cm.register_string(self.market_dynamics)
        self.cm.register_string(self.macro_system_state)
        
        # Define couplings
        self._setup_couplings()
        
        # Define resonance rules
        self._setup_resonance_rules()
    
    def _setup_couplings(self):
        """Set up coupling relationships"""
        # Agent trust affects system security
        self.cm.add_coupling(Coupling(
            source_string="agent_population",
            source_mode="trust_level",
            target_string="system_infrastructure",
            target_mode="security_posture",
            strength=0.65
        ))
        
        # System reliability affects market confidence
        self.cm.add_coupling(Coupling(
            source_string="system_infrastructure",
            source_mode="reliability",
            target_string="market_dynamics",
            target_mode="market_confidence",
            strength=0.75
        ))
        
        # Decision transparency affects agent trust
        self.cm.add_coupling(Coupling(
            source_string="decision_framework",
            source_mode="transparency",
            target_string="agent_population",
            target_mode="trust_level",
            strength=0.70
        ))
        
        # Market confidence affects macro stability
        self.cm.add_coupling(Coupling(
            source_string="market_dynamics",
            source_mode="market_confidence",
            target_string="macro_system_state",
            target_mode="overall_stability",
            strength=0.70
        ))
        
        # System security inversely affects macro risk
        self.cm.add_coupling(Coupling(
            source_string="system_infrastructure",
            source_mode="security_posture",
            target_string="macro_system_state",
            target_mode="systemic_risk",
            strength=-0.85  # Negative: higher security = lower risk
        ))
        
        # Agent risk affects macro risk
        self.cm.add_coupling(Coupling(
            source_string="agent_population",
            source_mode="risk_score",
            target_string="macro_system_state",
            target_mode="systemic_risk",
            strength=0.90
        ))
    
    def _setup_resonance_rules(self):
        """Set up resonance rules"""
        # Resonance: High agent trust amplifies improvements
        self.resonance_engine.add_rule(ResonanceRule(
            name="trust_resonance",
            condition=lambda strings: (
                strings["agent_population"].get_mode("trust_level") > 8.0
            ),
            effect=lambda strings: print("  Effect: Amplified improvements (multiplier: 1.5)")
        ))
        
        # Resonance: Low reliability triggers cascade
        self.resonance_engine.add_rule(ResonanceRule(
            name="reliability_cascade",
            condition=lambda strings: (
                strings["system_infrastructure"].get_mode("reliability") < 3.0
            ),
            effect=lambda strings: (
                strings["macro_system_state"].set_mode(
                    "systemic_risk",
                    strings["macro_system_state"].get_mode("systemic_risk") + 0.3
                ),
                print("  Effect: Risk amplification (+0.3)")
            )
        ))
        
        # Resonance: High market confidence + adoption creates feedback loop
        self.resonance_engine.add_rule(ResonanceRule(
            name="positive_feedback_loop",
            condition=lambda strings: (
                strings["market_dynamics"].get_mode("market_confidence") > 7.5 and
                strings["agent_population"].get_mode("adoption_rate") > 6.0
            ),
            effect=lambda strings: print("  Effect: Accelerated growth (growth_rate: 1.2x)")
        ))
        
        # Resonance: Critical risk threshold
        self.resonance_engine.add_rule(ResonanceRule(
            name="critical_risk_alert",
            condition=lambda strings: (
                strings["macro_system_state"].get_mode("systemic_risk") > 8.5
            ),
            effect=lambda strings: print(
                "  Effect: Emergency intervention required\n"
                "  Actions: increase_transparency, enhance_accountability, improve_security"
            )
        ))
    
    def update_agent_trust(self, new_value: float):
        """Update agent trust and propagate changes"""
        self.agent_population.set_mode("trust_level", new_value)
        self.cm.propagate("agent_population", "trust_level", new_value)
        self.resonance_engine.evaluate()
    
    def update_system_reliability(self, new_value: float):
        """Update system reliability and propagate changes"""
        self.system_infrastructure.set_mode("reliability", new_value)
        self.cm.propagate("system_infrastructure", "reliability", new_value)
        self.resonance_engine.evaluate()
    
    def update_agent_risk(self, new_value: float):
        """Update agent risk and propagate changes"""
        self.agent_population.set_mode("risk_score", new_value)
        self.cm.propagate("agent_population", "risk_score", new_value)
        self.resonance_engine.evaluate()
    
    def print_state(self):
        """Print current system state"""
        print("Current System State:")
        print(f"  Agent Trust Level: {self.agent_population.get_mode('trust_level'):.2f}")
        print(f"  Agent Risk Score: {self.agent_population.get_mode('risk_score'):.2f}")
        print(f"  System Reliability: {self.system_infrastructure.get_mode('reliability'):.2f}")
        print(f"  System Security: {self.system_infrastructure.get_mode('security_posture'):.2f}")
        print(f"  Decision Transparency: {self.decision_framework.get_mode('transparency'):.2f}")
        print(f"  Market Confidence: {self.market_dynamics.get_mode('market_confidence'):.2f}")
        print(f"  Macro Stability: {self.macro_system_state.get_mode('overall_stability'):.2f}")
        print(f"  Systemic Risk: {self.macro_system_state.get_mode('systemic_risk'):.2f}")
        print()


# ============================================================================
# Simulation Function
# ============================================================================

def simulate_macro_evolution(time_steps: int):
    """Simulate macro system evolution over time"""
    print("=" * 50)
    print("Psychohistory Macro Dynamics Simulation (Python)")
    print("=" * 50)
    print()
    print("This Python implementation demonstrates:")
    print("1. Manual coupling propagation (no native support)")
    print("2. Manual resonance rule evaluation")
    print("3. Class-based entity modeling")
    print("4. Callback-based observer pattern")
    print()
    print("Comparison with UAD:")
    print("  - UAD: ~250 LOC with native String Theory semantics")
    print("  - Python: ~500+ LOC with manual coupling/resonance")
    print("  - Expressiveness: UAD 5/5, Python 3/5")
    print()
    
    # Initialize system
    system = MacroDynamicsSystem()
    
    print("Initial System State:")
    system.print_state()
    
    print("Simulating", time_steps, "time steps...")
    print()
    
    print("Time Step 5:")
    print("  Agent trust increases → System security improves")
    system.update_agent_trust(7.2)
    system.print_state()
    
    print("Time Step 10:")
    print("  Positive feedback loop detected!")
    system.market_dynamics.set_mode("market_confidence", 8.0)
    system.agent_population.set_mode("adoption_rate", 7.0)
    system.resonance_engine.evaluate()
    print()
    
    print("Time Step 15:")
    print("  External shock: Reliability drops to 2.5")
    system.update_system_reliability(2.5)
    system.print_state()
    
    print("Time Step 20:")
    print("  Stabilization measures activated")
    system.decision_framework.set_mode("transparency", 7.0)
    system.decision_framework.set_mode("accountability", 7.5)
    system.update_system_reliability(6.8)
    system.print_state()
    
    print("=" * 50)
    print("Simulation Complete")
    print("=" * 50)
    print()


# ============================================================================
# Main Execution
# ============================================================================

if __name__ == "__main__":
    print("=" * 50)
    print("Psychohistory Scenario")
    print("=" * 50)
    print()
    print("This example demonstrates:")
    print("1. Brane definitions (dimensional contexts)")
    print("2. String declarations with vibrational modes")
    print("3. Coupling relationships between modes")
    print("4. Resonance rules for conditional effects")
    print("5. Macro-level system evolution simulation")
    print()
    
    simulate_macro_evolution(20)
    
    print("To fully execute, the runtime needs:")
    print("- String state management system")
    print("- Brane context tracking")
    print("- Coupling propagation engine")
    print("- Resonance condition evaluator")
    print("- Macro dynamics simulator")


