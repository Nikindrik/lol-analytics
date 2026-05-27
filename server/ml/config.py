from pathlib import Path

MODELS_DIR = Path('models')

STYLE_NAMES = {
    'ADC': {0: 'Consistent', 1: 'Careful', 2: 'Aggressive'},
    'Jungle': {0: 'Standard', 1: 'Passive Ganker', 2: 'Aggressive Farm'},
    'Mid': {0: 'Balanced', 1: 'Passive', 2: 'Aggressive'},
    'Top': {0: 'Standard', 1: 'Passive', 2: 'Aggressive'},
    'Support': {
        0: 'Safe Roaming',
        1: 'Consistent',
        2: 'Passive',
        3: 'Specialist',
        4: 'Specialist',
        5: 'Specialist',
        6: 'Specialist'
    }
}
