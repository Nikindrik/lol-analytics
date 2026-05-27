from config import MODELS_DIR, STYLE_NAMES
from ensemble_model import EnsembleModel


class MLService:
    def __init__(self):
        self.models = {}
        for role_dir in MODELS_DIR.iterdir():
            if role_dir.is_dir():
                self.models[role_dir.name] = EnsembleModel(role_dir)
                print(f"Loaded model for {role_dir.name}")

    def predict(self, role, data):
        print(f'Predict called with role={role}, data={data}')

        role = role or 'Support'
        if role not in self.models:
            print(f'Role {role} not found, using Support')
            role = 'Support'

        cluster, confidence = self.models[role].predict(data)
        style = STYLE_NAMES.get(role, {}).get(cluster, f'Style {cluster}')

        print(f'Prediction result: cluster={cluster}, style={style}, confidence={confidence}')

        return {
            'cluster': int(cluster),
            'style': style,
            'confidence': confidence
        }
