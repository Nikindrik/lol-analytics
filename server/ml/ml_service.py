from config import MODELS_DIR, STYLE_NAMES
from ensemble_model import EnsembleModel
from regression_ensemble import RegressionEnsemble  # НОВЫЙ импорт



class MLService:
    def __init__(self):
        self.models = {}
        self.regression_models = {}  # НОВОЕ: словарь для регрессионных моделей
        for role_dir in MODELS_DIR.iterdir():
            if role_dir.is_dir():
                self.models[role_dir.name] = EnsembleModel(role_dir)
                print(f"Loaded model for {role_dir.name}")
                if self._has_regression_files(role_dir):
                    self.regression_models[role_dir.name] = RegressionEnsemble(role_dir)
                    print(f"Loaded regression model for {role_dir.name}")

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
    def _has_regression_files(self, role_dir):
        """НОВЫЙ метод: проверяет наличие файлов для регрессии"""
        # Проверяем наличие файлов моделей, которые могут быть для регрессии
        # (используем те же файлы что и для классификации, но загружаем как регрессию)
        has_model = False
        for model_file in ['xgboost.json', 'lightgbm.txt', 'catboost.cbm']:
            if (role_dir / model_file).exists():
                has_model = True
                break
        
        # Также нужен features.json
        has_features = (role_dir / 'features.json').exists()
        
        return has_model and has_features
        # НОВЫЙ метод для регрессии
    def predict_regression(self, role, data):
        """Предсказание регрессии (например, золота) для конкретной роли"""
        role = role or 'Support'
        
        # Если нет регрессионной модели для этой роли, возвращаем None
        if role not in self.regression_models:
            return None
        
        # Извлекаем признаки в правильном порядке (как в feature_extractor)
        from feature_extractor import FeatureExtractor
        extractor = FeatureExtractor(self.regression_models[role].features)
        X_raw = extractor.transform(data)
        
        # Предсказываем
        prediction = self.regression_models[role].predict(X_raw)
        
        return prediction