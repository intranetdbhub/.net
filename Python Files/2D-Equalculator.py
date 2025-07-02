#Imports
from tkinter import *
from cmath import sqrt # cmath is used instead of math to avoid a math domain error when trying to get the square root of a negative number in the discriminant sec

# Creating the window
window = Tk()
window.geometry("700x385")  # width x height standard measures of the GUI window
window.title("2D-Equalculator")
window.configure(background="peach puff")

# Main Title
Label(window, text="         Quadratic Equation Calculator (x^2+bx+c=0)", bg="peach puff", fg="black", font="Helvetica 18 bold").grid(row=0,
                                                                                                          column=0,
                                                                                                          sticky=N)

# User's input box labels and text entry boxes for coeficient (a):
Label(window, text="                              Type coeficient (a) there ---> ", bg="peach puff",
      fg="black", font="Helvetica 14 italic").grid(row=11, column=0, sticky=W)
textentry1 = Entry(window, width=7, bg="white")
textentry1.grid(row=11, column=1, sticky=W)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=12, column=0, sticky=W)


# User's input box labels and text entry boxes for coeficient (b):
Label(window, text="                              Type coeficient (b) there ---> ", bg="peach puff",
      fg="black", font="Helvetica 14 italic").grid(row=13, column=0, sticky=W)
textentry2 = Entry(window, width=7, bg="white")
textentry2.grid(row=13, column=1, sticky=W)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=14, column=0,
                                                                                           sticky=W)

# User's input box labels and text entry boxes for coeficient (c):
Label(window, text="                              Type coeficient (c) there ---> ", bg="peach puff",
      fg="black", font="Helvetica 14 italic").grid(row=15, column=0, sticky=W)
textentry3= Entry(window, width=7, bg="white")
textentry3.grid(row=15, column=1, sticky=W)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=16, column=0,
                                                                                           sticky=W)

# Creating an output text box to display results. It needs to be resolved below to display two results inside one box
output = Text(window, width=45, height=3, wrap=WORD, background="white")
output.grid(row=20, column=0, columnspan=3, sticky=N)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=21, column=1,sticky=W)

#Second Degree Ecuations Calculations Function when clicking the Calculate Roots button
def Calculate_Equation_Roots():

    # This will collect text from the entry boxes
    entered_text1 = float(textentry1.get()) #coeficient a
    entered_text2 = float(textentry2.get()) #coeficient b
    entered_text3 = float(textentry3.get()) #coeficient c

    #Calculating the discriminant (y)
    pre_y = (entered_text2*entered_text2) - (4*entered_text1*entered_text3) #a(entered_text1), b(entered_text2), c(entered_text3)
    y = sqrt(pre_y)

    # Calculating the roots (x,y):
    x1 = ((-1 * entered_text2) + y) / (2 * entered_text1)
    x2 = ((-1 * entered_text2) - y) / (2 * entered_text1)

    # The return that will print results in the output box:
    x1 = "The value of the positive & negative roots are: {0}  ".format(x1)
    output.insert(END, x1)
    x2 = "                                                     {0}  ".format(x2)
    output.insert(END, x2)


# Clear all button function when clicked
def clear_textentries():
    textentry1.delete(0, END)
    textentry2.delete(0, END)
    textentry3.delete(0, END)
    output.delete("1.0", "end") #Different from the textentries boxes you put "1.0", "end"

# Calculate roots button creation and its function declaration
Button(window, text="Calculate Roots", width=12, command=Calculate_Equation_Roots).grid(row=17, column=1, sticky=W)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=18, column=0,
                                                                                               sticky=W)
# Clear All button creation and clear_textentries function declaration
Button(window, text="Clear All", width=9, command=clear_textentries).grid(row=26, column=1, sticky=E)
Label(window, text=" ", bg="peach puff", fg="peach puff", font="Helvetica 14 italic").grid(row=16, column=0,
                                                                                           sticky=W)


window.mainloop()
