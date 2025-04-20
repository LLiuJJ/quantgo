# Quantgo
A Quantitative Trading System.

# Example

1.Query the stock price of Apple Inc. in the past 100 days and plot it。

```
 ./quantgo-cli plot_stock_chart -c AAPL -d 100
```

![apple stock price history 100 days](assets/newplot.png)

2.Quantitative investment decision-making of Alibaba stock using turtle trading algorithm。

```
./quantgo-cli quant_cal -a turtle -c BABA -d 50
```

![BABA stock price history 50 days,buy decision](assets/decision.png)


![BABA stock price portfolio using turtle trading algorithm](assets/portfolio.png)

