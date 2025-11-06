import marimo

__generated_with = "0.17.2"
app = marimo.App(width="medium")


@app.cell
def _():
    import marimo as mo
    import polars as pl
    import matplotlib.pyplot as plt

    __generated_with__ = "0.17.2"
    app = mo.App(width="medium")
    return mo, pl, plt


@app.cell
def _(mo):
    mo.md(
        """
    # Week 2 Assignment

    #### **Dataset name:** modified_diamonds.csv  
    #### **Due Date:** Wednesday, November 5th, 6:00 pm mountain time

    As demonstrated in the class notebook for week two and by what we covered in class, you will provide an analysis for the modified_diamonds dataset. Assume the dataset represents the current regional inventories for a company that sells diamonds.

    **Assignment requirements:**
    1. Write a data dictionary for the modified_diamonds dataset (CSV file).  
       - column_name, description, data_class, data_sub_class, datatype, unit_of_measure  
    2. Ensure correct data types for each column; verify via `.schema`.  
    3. Develop a high-level question or business scenario and analyze it.  
    4. Define four univariate analytical questions (two segmented).  
    5. Use descriptive statistics and at least one visualization for each question.  
    6. Discuss what you are doing and why.  
    7. Recap insights and provide a summary statement answering the scenario.
    """
    )
    return


@app.cell
def _(mo):
    mo.md(
        """
    ## Project Week 2 Check-in
    Create **one** notebook for all check-ins and final submission.

    **Requirements for Week 2 Check-in:**
    1. Create a data dictionary in CSV format.  
    2. Import and verify datatypes.  
    3. Develop a high-level question or scenario and discuss it.  
    4. Define and analyze 3 analytical questions relevant to your scenario.
    """
    )
    return


@app.cell
def _(mo):
    mo.md(
        """
    ### **Deliverables**
    1. The text file with the dataset data dictionary.  
    2. This Marimo notebook with all requirements completed.  
    3. The course project notebook.

    **Grading:** 50% based on analysis (statistics/visuals) and 50% on discussion/narrative.
    """
    )
    return


@app.cell
def _(pl):
    #Import dataset:
    df = pl.read_csv("modified_diamonds.csv", infer_schema_length=50000)
    return (df,)


@app.cell
def _(df, mo):
    # View the schema for the data dictionary:
    mo.md(f"**Rows:** {df.height} **Cols:** {len(df.columns)}")
    mo.md(f"**Schema:** `{df.schema}`")
    return


@app.cell
def _(df, mo, pl):
    # Note: All mathematical/statistical formulas in the entire code below was written with the assitance of ChatGPT and Claude AI models:

    # Data dictionary creation for the modified_diamonds dataset:
    desc = {
            "carat": "Diamond weight in carats.",
            "cut": "Cut quality grade (Ideal, Very Good, Premium, Good, Fair).",
            "cut_rank": "Ordinal rank of cut quality (1 = best).",
            "color": "Color grade (D best/colourless → J more tinted).",
            "color_rank": "Ordinal rank of color (1 = best).",
            "clarity": "Clarity grade (IF, VVS1/2, VS1/2, SI1/2, I1).",
            "clarity_rank": "Ordinal rank of clarity (for example: 1 = best).",
            "depth_pct": "Total depth as percent of average diameter.",
            "table_mm": "Table (top facet) width in millimeters.",
            "price_usd": "Price in US dollars.",
            "length_mm": "Length (longest dimension) in millimeters.",
            "width_mm": "Width (shortest face-to-face) in millimeters.",
            "depth_mm": "Depth (table to culet) in millimeters.",
            "region": "Inventory region / warehouse label.",
        }

    units = {
        "carat": "carat",
        "depth_pct": "%",
        "table_mm": "mm",
        "price_usd": "USD",
        "length_mm": "mm",
        "width_mm": "mm",
        "depth_mm": "mm",
        }

    ordinal_cats = {"cut", "color", "clarity"}
    ordinal_rank_nums = {"cut_rank", "color_rank", "clarity_rank"}

    rows = []
    for name, dtype in df.schema.items():
        dt = str(dtype)
        is_int = "Int" in dt
        is_float = "Float" in dt
        is_num = is_int or is_float

        if name in ordinal_cats:
           data_class, data_sub = "categorical", "ordinal"
        elif name in ordinal_rank_nums:
           data_class, data_sub = "numerical", "discrete"
        else:
           data_class = "numerical" if is_num else "categorical"
           data_sub = ("discrete" if is_int else "continuous") if is_num else "nominal"

        rows.append({
            "column_name": name,
            "description": desc.get(name, ""),
            "data_class": data_class,
            "data_sub_class": data_sub,
            "datatype": dt,
            "unit_of_measure": units.get(name, "")
            })

    dd = pl.DataFrame(rows)
    dd.write_csv("modified_diamonds_data_dictionary.csv")
    mo.md("*Data dictionary saved to `modified_diamonds_data_dictionary.csv`*")
    dd
    return


@app.cell
def _(mo):
    mo.md(
        """
    ## High Level Business Scenario:

    The company Intranet Innovations And Services LLC, manages regional inventories of diamonds. Management wants to optimize pricing and regional distribution by understanding which and how the above product attributes influence price. They also want to know which regions hold the highest-value stock.

    ### Brief Analysis
    The above modified_diamonds_data_dictionary.csv dictionary file provides additional details concerning product attributes. as you can see, each diamond record includes physical characteristics (carat, dimensions, clarity, color, cut) and a price in US dollars. The modified_diamonds.csv dataset shown in the above code from which the modified_diamonds_data_dictionary.csv dictionary file comes from,  represents current diamond inventory from multiple regions (region1–region4). As a business analyst I was able to use python and descriptive analytics to identify which attributes most affect price. Also, I determined whether certain regions hold more valuable or higher-quality diamonds, by answering 3 analytical questions. Lastly, I will show below through data visualization, how do diamond characteristics and regional distribution influence pricing and inventory value for Intranet Innovations And Services LLC.
    """
    )
    return


@app.cell
def _(mo):
    mo.md(
        """
    ## Analytical Question #1:
    What is the relationship between diamond size (carat) and price?

    ### Analytical Question #1 Analysis:
    Using univariate analysis, the scatter or histogram shown in below cell of carat vs price_usd, shows that when diamond size (carat weight) increases, prices generally trend upward, with smaller diamonds under 0.6 carats being more affordable and larger stones above 1.0 carat having higher prices. However, while carat is a strong price driver, other clarity, color, and cut also significantly influence pricing. For Intranet Innovations And Services LLC, this means that carat weight should be the primary pricing factor, but all diamond attributes must be evaluated together. The company's inventory should focus on mid-size diamonds (0.7-1.2 carats) to help balance customer affordability with strong profit margins. See the code below and execute it to view the data visualization:
    """
    )
    return


@app.cell
def _(df, pl, plt):
    # Analytical Question #1 Data Visualization — Average Price by Carat Bin:
    bin_width = 0.2
    binned = (
        df.with_columns((pl.col("carat") / bin_width).floor() * bin_width)
          .rename({"carat": "carat_binned"})
    )

    agg = (
        binned.group_by("carat_binned")
              .agg(pl.col("price_usd").mean().alias("avg_price"))
              .sort("carat_binned")
    )

    plt.figure(figsize=(8, 4))
    plt.bar(
        [f"{v:.2f}" for v in agg["carat_binned"].to_list()],
        agg["avg_price"].to_list(),
        color="steelblue"
    )
    plt.xlabel("Carat (binned)")
    plt.ylabel("Average Price (USD)")
    plt.title("Average Diamond Price by Carat Bin")
    plt.xticks(rotation=45, ha="right")
    plt.show()
    return


@app.cell
def _(mo):
    mo.md(
        """
    ## Analytical Question #2:
    What is the correlation between diamond cut quality and higher pricing?

    ### Analytical Question #2 Analysis:
    That better cut grades don't automatically mean higher prices. Ideal cuts (the best quality) showed the lowest average price at around the 3,900s USD, while "Premium" cuts have the highest at around $5,000. This is because a diamond's price depends on multiple factors working together and not just the cut. Premium cut diamonds tend to be found as larger or have better clarity, which makes them cost more. Also, Ideal cuts may be more common in smaller diamonds.
    Intranet Innovations And Services LLC, should emphasize sparkle and craftsmanship, and not use cut grades to automatically predict price on its own. Also, the focus should be more on size and clarity, abnd use cut quality as a selling point for beauty rather than as a primary price driver. See the code below and execute it to view the data visualization:
    """
    )
    return


@app.cell
def _(df, pl, plt):
    # Analytical Question #2 Data Visualization: Average Price by Cut Quality:
    cut_stats = (
        df.group_by("cut")
          .agg(pl.col("price_usd").mean().alias("avg_price"))
          .sort("avg_price")
    )

    # Define a distinct color for each bar
    colors = ["#FF5733", "#FFC300", "#DAF7A6", "#33FFCE", "#3375FF"]

    plt.figure(figsize=(8, 4))
    plt.bar(cut_stats["cut"].to_list(), cut_stats["avg_price"].to_list(), color=colors)
    plt.xlabel("Cut Quality")
    plt.ylabel("Average Price (USD)")
    plt.title("Average Diamond Price by Cut Quality")
    plt.show()
    return


@app.cell
def _(mo):
    mo.md(
        """
    ## Analytical Question #3:
    Do regions differ in their average diamond prices?

    ### Analytical Question #3 Analysis:
    Yes. The average diamond prices vary by region, with region1 having the highest average price around the 4,900s USDs, followed by region2 around the 4600s USDs, region4 around the 4,200s USDs, and region3 around the 3,600s. But, we see that the price differences are modest between highest and lowest. This suggests that regional variation is a minor factor compared to diamond characteristics like carat, clarity, and color. The higher prices in region1 and region2 likely reflect customer preferences or purchasing power in those markets rather than fundamentally different inventory needs. Intranet Innovations And Services LLC, needs to maintain similar inventory quality across all regions while potentially adjusting the mix of diamond sizes and price points to match local demand, offering more premium options in region1 and more budget-friendly selections in region3. Execute the code below so the data visualization can be exposed:
    """
    )
    return


@app.cell
def _(df, pl, plt):
    # Analytical Question#3 Data Visualization:
    region_stats = (
        df.group_by("region")
          .agg(pl.col("price_usd").mean().alias("avg_price"))
          .sort("avg_price")
    )

    # Unique colors for each region
    region_colors = ["#FF9999", "#66B2FF", "#99FF99", "#FFCC99"]

    plt.figure(figsize=(7, 4))
    plt.bar(
        region_stats["region"].to_list(),
        region_stats["avg_price"].to_list(),
        color=region_colors
    )
    plt.xlabel("Region")
    plt.ylabel("Average Price (USD)")
    plt.title("Average Diamond Price by Region")
    plt.show()
    return


if __name__ == "__main__":
    app.run()
