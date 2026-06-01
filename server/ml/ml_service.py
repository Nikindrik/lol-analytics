from config import MODELS_DIR, STYLE_NAMES
from ensemble_model import EnsembleModel
from regression_ensemble import RegressionEnsemble


class MLService:
    def __init__(self):
        self.models = {}
        self.regression_models = {}
        
        for role_dir in MODELS_DIR.iterdir():
            if role_dir.is_dir():
                self.models[role_dir.name] = EnsembleModel(role_dir)
                print(f"Loaded model for {role_dir.name}")
                
                # Загружаем регрессионную модель из подпапки gold_prediction
                gold_model_path = role_dir / 'gold_prediction'
                if gold_model_path.exists() and self._has_regression_files(gold_model_path):
                    self.regression_models[role_dir.name] = RegressionEnsemble(gold_model_path)
                    print(f"Loaded regression model for {role_dir.name} from gold_prediction")

    def _has_regression_files(self, model_path):
        """Проверяет наличие файлов для регрессии"""
        has_model = False
        for model_file in ['xgboost.json', 'lightgbm.txt', 'catboost.cbm']:
            if (model_path / model_file).exists():
                has_model = True
                break
        
        has_features = (model_path / 'features.json').exists()
        
        return has_model and has_features

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
    
    def predict_regression(self, role, data):
        """Предсказание регрессии (золота) для конкретной роли"""
        role = role or 'Support'
        
        if role not in self.regression_models:
            print(f"No regression model for role {role}")
            return None
        
        try:
            from feature_extractor import FeatureExtractor
            extractor = FeatureExtractor(self.regression_models[role].features)
            X_raw = extractor.transform(data)
            prediction = self.regression_models[role].predict(X_raw)
            print(f"Regression prediction for {role}: {prediction}")
            return prediction
        except Exception as e:
            print(f"Regression error: {e}")
            return None
    
    def predict_full(self, role, data):
        """Полное предсказание: стиль + золото"""
        playstyle_result = self.predict(role, data)
        gold_value = self.predict_regression(role, data)
        
        result = {'playstyle': playstyle_result}
        
        if gold_value is not None:
            current_gold = data.get('current_gold', data.get('total_gold', 0))
            result['gold_prediction'] = {
                'predicted_gold': round(float(gold_value), 0),
                'current_gold': round(float(current_gold), 0),
                'gold_diff': round(float(gold_value) - float(current_gold), 0)
            }
        
        return result