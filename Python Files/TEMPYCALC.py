from tkinter import *

# Creating the window
window = Tk()
window.geometry("630x560")  # width x height standard measures of the GUI window
window.title("TEMPYCALC")
window.configure(background="peach puff")

# Main Title
Label(window, text="Temperature Calculator ", bg="peach puff", fg="black", font="Helvetica 18 bold").grid(row=0,
                                                                                                          column=0,
                                                                                                          sticky=N)
# Menu's option labels
Label(window, text="Choose one of the following options: ", bg="peach puff", fg="black", font="Helvetica 16").grid(
    row=3, column=0, sticky=W)
Label(window, text="1-) Convert from Celsius to Fahrenheit. ", bg="peach puff", fg="black",
      font="Helvetica 14 italic").grid(row=4, column=0, sticky=W)
Label(window, text="2-) Convert from Fahrenheit to Celsius. ", bg="peach puff", fg="black",
      font="Helvetica 14 italic").grid(row=5, column=0, sticky=W)
Label(window, text="3-) Convert from Celsius to Kelvin. ", bg="peach puff", fg="black",
      font="Helvetica 14 italic").grid(row=6, column=0, sticky=W)
Label(window, text="4-) Convert from Fahrenheit to Kelvin. ", bg="peach puff", fg="black",
      font="Helvetica 14 italic").grid(row=7, column=0, sticky=W)
Label(window, text="5-) Convert from Kelvin to Celsius. ", bg="peach puff", fg="black",
      font="Helvetica 14 italic").grid(row=8, column=0, sticky=W)
Label(window, text="6-) Convert from Kelvin to Fahrenheit. ", bg="peach puff", fg="black",
      font="Helvetica 14 italic").grid(row=9, column=0, sticky=W)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=10, column=0,
                                                                                           sticky=W)  # space

# User's input box 1
Label(window, text="                              Type the number of your chosen option here: ", bg="peach puff",
      fg="black", font="Helvetica 14 italic").grid(row=11, column=0, sticky=W)
textentry1 = Entry(window, width=4, bg="white")
textentry1.grid(row=11, column=1, sticky=W)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=12, column=0, sticky=W)


# User's input box 2
Label(window, text="                          Type the temperature you want to convert here: ", bg="peach puff",
      fg="black", font="Helvetica 14 italic").grid(row=13, column=0, sticky=W)
textentry2 = Entry(window, width=7, bg="white")
textentry2.grid(row=13, column=1, sticky=W)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=14, column=0,
                                                                                           sticky=W)  # space


# Creating an output text box
output = Text(window, width=45, height=3, wrap=WORD, background="white")
output.grid(row=17, column=0, columnspan=3, sticky=N)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=18, column=0,sticky=W)


# Convert button when clicked
def click():
    # This will collect text from the entry boxes
    entered_text1 = float(textentry1.get())
    entered_text2 = float(textentry2.get())

    # Formulas to do the calculations
    f1 = ((entered_text2 * (9 / 5)) + 32)
    f2 = (entered_text2 - 32) * (5 / 9)
    f3 = entered_text2 + 273.15
    f4 = ((entered_text2 - 32) * (5 / 9)) + 273.15
    f5 = entered_text2 - 273.15
    f6 = ((entered_text2 - 273.15) * (9 / 5)) + 32

# Calculating according to options
    if entered_text1 == 1:
        tf1 = "Your Celsius temperature converted into Fahrenheit is: {0}".format(f1)
        output.insert(END, tf1)  # The return that will print results in the output box
    elif entered_text1 == 2:
        tf2 = "Your Fahrenheit temperature converted into Celsius is: {0}".format(f2)
        output.insert(END, tf2)  # The return that will print results in the output box
    elif entered_text1 == 3:
        tf3 = "Your Celsius temperature converted into Kelvin is: {0}".format(f3)
        output.insert(END, tf3)  # The return that will print results in the output box
    elif entered_text1 == 4:
        tf4 = "Your Fahrenheit temperature converted into Kelvin is: {0}".format(f4)
        output.insert(END, tf4)  # The return that will print results in the output box
    elif entered_text1 == 5:
        tf5 = "Your Kelvin temperature converted into Celsius is: {0}".format(f5)
        output.insert(END, tf5)  # The return that will print results in the output box
    elif entered_text1 == 6:
        tf6 = "Your Kelvin temperature converted into Fahrenheit is: {0}".format(f6)
        output.insert(END, tf6)  # The return that will print results in the output box
    else:
        incorrect_option = "You chose the incorrect option or you have an error in your temperature's input!!!"
        output.insert(END, incorrect_option)  # The return that will print results in the output box


# Clear all button when clicked
def clear_textentries():
    textentry1.delete(0, END)
    textentry2.delete(0, END)
    output.delete("1.0", "end") #Different from the textentries boxes you put "1.0", "end"

# Convert button creation and click function declaration
Button(window, text="Convert", width=9, command=click).grid(row=15, column=1, sticky=W)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=16, column=0,
                                                                                       sticky=W)

# Clear All button creation and clear_textentries function declaration
Button(window, text="Clear All", width=9, command=clear_textentries).grid(row=19, column=1, sticky=E)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=16, column=0,
                                                                                           sticky=W)


# main root loop run
window.mainloop()
