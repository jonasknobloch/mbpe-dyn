package main

var types = [3]string{
	"target",
	"base",
	"instruction",
}

var instructions = [3][4]string{
	{
		"Nominalized adjective:",
		"Noun:",
		"The following is a nominalized adjective:",
		"The following is a noun:",
	},
	{
		"{} ->",
		"{} :",
		"{} -",
		"{}",
	},
	{
		"Adjective: {}\nNominalization:",
		"Form the nominalization of the given adjective.\n\n{} ->",
		"Nominalize the given adjective.\n\n{} ->",
		"Turn the given adjective into a noun.\n\n{} ->",
	},
}
