import json

import catboost
import joblib
import lightgbm as lgb
import xgboost as xgb


class ModelLoader:
    def __init__(self, role_dir):
        self.role_dir = role_dir

    def load_catboost(self):
        path = self.role_dir / 'catboost.cbm'
        if not path.exists():
            return None
        model = catboost.CatBoostClassifier()
        model.load_model(str(path), format='cbm')
        return model

    def load_xgboost(self):
        path = self.role_dir / 'xgboost.json'
        if not path.exists():
            return None
        model = xgb.XGBClassifier()
        model.load_model(str(path))
        return model

    def load_lightgbm(self):
        path = self.role_dir / 'lightgbm.txt'
        if not path.exists():
            return None
        return lgb.Booster(model_file=str(path))

    def load_scaler(self):
        path = self.role_dir / 'scaler.pkl'
        if not path.exists():
            return None
        return joblib.load(path)

    def load_features(self):
        path = self.role_dir / 'features.json'
        if not path.exists():
            return []
        with open(path) as f:
            return json.load(f)
