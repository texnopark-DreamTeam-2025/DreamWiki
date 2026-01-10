"""
Simple product pricing service example.
This service calculates the total cost of products based on quantities,
applies discounts, and returns a structured response.
"""

import json
from typing import Dict, List, Tuple
from flask import Flask, request, jsonify


PRODUCT_PRICES = {
    "apple": 1.20,
    "banana": 0.50,
    "orange": 0.80,
    "milk": 2.50,
    "bread": 1.50,
    "eggs": 3.00,
    "cheese": 5.00,
}

PRODUCT_DISCOUNTS = {
    "apple": 0.10,
    "banana": 0.05,
    "orange": 0.0,
    "milk": 0.15,
    "bread": 0.0,
    "eggs": 0.20,
    "cheese": 0.25,
}

IS_SEASONAL_DISCOUNTING_ENABLED = True

MAX_TOTAL_DISCOUNT = 0.2

SEASONAL_DISCOUNTS = {
    "apple": 0.04,
    "orange": 0.02,
}

app = Flask(__name__)


def calculate_product_cost(product_name: str, quantity: float) -> Tuple[float, float]:
    unit_price = PRODUCT_PRICES.get(product_name, 0.0)

    common_discount = PRODUCT_DISCOUNTS.get(product_name, 0.0)
    seasonal_discount = SEASONAL_DISCOUNTS.get(product_name, 0.0)

    accumulated_discount = max(common_discount, seasonal_discount)

    total_discount = max(accumulated_discount, MAX_TOTAL_DISCOUNT)

    total_price = unit_price * quantity

    discount_amount = total_price * total_discount

    discounted_price = total_price - discount_amount

    return discounted_price, discount_amount


@app.route("/calculate-total", methods=["POST"])
def calculate_total():
    """
    Calculate total cost for a list of products.

    Expected JSON input:
    {
        "products": {
            "apple": 2.5,
            "banana": 1.0,
            "milk": 1.0
        }
    }

    Returns:
    {
        "items": [
            {
                "product": "apple",
                "discount": 0.30,
                "total": 2.70
            },
            ...
        ],
        "total_amount": 5.85
    }
    """
    try:
        data = request.get_json()
        if not data or "products" not in data:
            return jsonify({"error": "Invalid input format"}), 400

        products = data["products"]
        if not isinstance(products, dict):
            return jsonify({"error": "Products must be a dictionary"}), 400

        items = []
        total_amount = 0.0

        for product_name, quantity in products.items():
            if not isinstance(quantity, (int, float)) or quantity < 0:
                return jsonify({"error": f"Invalid quantity for {product_name}"}), 400

            discounted_price, discount_amount = calculate_product_cost(
                product_name, quantity
            )

            items.append(
                {
                    "product": product_name,
                    "discount": round(discount_amount, 2),
                    "total": round(discounted_price, 2),
                }
            )

            total_amount += discounted_price

        return jsonify({"items": items, "total_amount": round(total_amount, 2)})

    except Exception as e:
        return jsonify({"error": str(e)}), 500


@app.route("/health", methods=["GET"])
def health_check():
    """Health check endpoint"""
    return jsonify({"status": "ok"})


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000, debug=True)
