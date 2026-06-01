from flask import Flask, jsonify, request
from ml_service import MLService

app = Flask(__name__)
service = MLService()

@app.route('/predict', methods=['POST'])
def predict():
    try:
        data = request.json
        role = data.get('role')

        features = {
            'kills': data.get('kills', 0),
            'deaths': data.get('deaths', 0),
            'assists': data.get('assists', 0),
            'items_count': data.get('items_count', 0),
            'cs': data.get('cs', 0),
            'jungle_cs': data.get('jungle_cs', 0),
            'gold_diff': data.get('gold_diff', 0),
            'xp_diff': data.get('xp_diff', 0),
            'team_gold_diff': data.get('team_gold_diff', 0),
            'jungle_ratio': data.get('jungle_ratio', 0),
            'kill_participation': data.get('kill_participation', 0),
        }

        result = service.predict(role, features)
        return jsonify(result)

    except Exception as e:
        return jsonify({'error': str(e)}), 500

@app.route('/health', methods=['GET'])
def health():
    return jsonify({
        'status': 'ok',
        'models': list(service.models.keys())
    })

if __name__ == '__main__':
    app.run(host='127.0.0.1', port=5000, debug=False)
    
@app.route('/predict_full', methods=['POST'])
def predict_full():
    try:
        data = request.json
        role = data.get('role')
        
        features = {
            'kills': data.get('kills', 0),
            'deaths': data.get('deaths', 0),
            'assists': data.get('assists', 0),
            'items_count': data.get('items_count', 0),
            'cs': data.get('cs', 0),
            'jungle_cs': data.get('jungle_cs', 0),
            'gold_diff': data.get('gold_diff', 0),
            'xp_diff': data.get('xp_diff', 0),
            'team_gold_diff': data.get('team_gold_diff', 0),
            'jungle_ratio': data.get('jungle_ratio', 0),
            'kill_participation': data.get('kill_participation', 0),
            'current_gold': data.get('current_gold', data.get('total_gold', 0)),
            'total_gold': data.get('total_gold', 0)
        }
        
        result = service.predict_full(role, features)
        return jsonify(result)
    
    except Exception as e:
        return jsonify({'error': str(e)}), 500
@app.route('/predict_regression', methods=['POST'])
def predict_regression():
    try:
        data = request.json
        role = data.get('role')
        
        features = {
            'kills': data.get('kills', 0),
            'deaths': data.get('deaths', 0),
            'assists': data.get('assists', 0),
            'items_count': data.get('items_count', 0),
            'cs': data.get('cs', 0),
            'jungle_cs': data.get('jungle_cs', 0),
            'gold_diff': data.get('gold_diff', 0),
            'xp_diff': data.get('xp_diff', 0),
            'team_gold_diff': data.get('team_gold_diff', 0),
            'jungle_ratio': data.get('jungle_ratio', 0),
            'kill_participation': data.get('kill_participation', 0),
        }
        
        prediction = service.predict_regression(role, features)
        
        if prediction is None:
            return jsonify({'error': f'No regression model for role {role}'}), 404
        
        return jsonify({
            'role': role,
            'predicted_value': prediction
        })
    
    except Exception as e:
        return jsonify({'error': str(e)}), 500

if __name__ == '__main__':
    app.run(host='127.0.0.1', port=5000, debug=False)