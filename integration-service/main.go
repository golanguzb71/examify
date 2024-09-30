package main

import (
	"context"
	"fmt"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
)

func main() {
	ctx := context.Background()

	apiKey := "AIzaSyBeQQyXZL0Duo-K36pDbTRM4EDi6thAMjo"

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	model.SetTemperature(0)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "application/json"

	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("The pie charts show the amount of revenue and expenditures in 2016 for a children’s charity in the USA. Overall, it can be seen that donated food accounted for the majority of the income, while program services accounted for the most expenditure. Total revenue sources just exceeded outgoings. In detail, donated food provided most of the revenue for the charity, at 86%. Similarly, with regard to expenditures, one category, program services, accounted for nearly all of the outgoings, at 95.8%. The other categories were much smaller. Community contributions, which were the second largest revenue source, brought in 10.4% of overall income, and this was followed by program revenue, at 2.2%. Investment income, government grants, and other income were very small sources of revenue, accounting for only 0.8% combined. There were only two other expenditure items, fundraising and management and general, accounting for 2.6% and 1.6% respectively. The total amount of income was $53,561,580, which was just enough to cover the expenditures of $53,224,896.\n\ngive answer like that \n```json\n{\n{\n  \"feedback\": \"\",\n  \"coherence_score\": this_task_coherence_score,\n  \"grammar_score\": \"grammar\",\n  \"lexical_resource_score\": lexical_resource_score,\n  \"task_achievement_score\": task_achievement_score,\n  \"task_band_score\": task_band_score\n}\n}\n```\nfor this \n\ngive answer like ielts coherense score 8 grammar score 7 like that\n"),
			},
		},
		{
			Role: "model",
			Parts: []genai.Part{
				genai.Text("```json\n{\n  \"feedback\": \"The provided text is a well-written summary of the pie charts, presenting a clear and concise overview of the charity's revenue and expenditure. The information is organized logically and presented in a way that is easy to understand.  The use of percentages and the specific dollar amounts for total revenue and expenditure further enhance the clarity of the information.\",\n  \"coherence_score\": 6,\n  \"grammar_score\": 6,\n  \"lexical_resource_score\": 6,\n  \"task_achievement_score\": 6,\n  \"task_band_score\": 6 \n}\n``` \n"),
			},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text("An increasing number of married couples around the world choosing to remain childless. The main benefits of not having a child for couples are that they can focus on their careers and have more time for themselves. The main drawbacks are that they could not fit into their peers’ group and have no one to look after them when they get old. \n\nOne primary advantage of remaining childless for married couples is that they can focus on their work. This is because they have less responsibility and distractions in their lives compared to the couples that have a child. Another advantage of this is that they have more spare time. Looking after a child is a full-time job for parents and taking most of their time, while child-free couples have lots of free time after work. For example, many couples stop going out late with their friends after having a child as they have to stay at home for looking after their children. \n\nOne disadvantage of couples deciding not to have children is that they can struggle to hang with their peers after most of them have children. Most parents prefer to spend more time with other couples that have children as well. Moreover, do not have anyone to look after them in their elderliness is another disadvantage. Children are the ones who take care of their parents when they get old because their parents did the same for them when they were young. For instance, the vast majority of the people who live in care homes have no child. \n\nIn conclusion, the main benefits of staying child-free for couples are that they can be more career-oriented and have more free time for themselves, and the main drawbacks are that they could have problems about fitting into their friends’ group and having no one to take care of them when they become older."))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		fmt.Printf("%v\n", part)
	}
}
