package bot

import (
	"fmt"
	"os"

	"github.com/0xfyer/discord-bot/storage"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
	state   *storage.State
}

func New(secret string) (*Bot, error) {
	session, err := discordgo.New("Bot " + secret)
	if err != nil {
		return nil, err
	}

	return &Bot{
		session: session,
		state:   storage.New(),
	}, nil
}

// Open() creates the bot session with Discord and
// blocks until the stop channel receives an os.Signal.
// You should almost always call this from a goroutine.
func (b *Bot) Open(stop chan os.Signal) error {
	err := b.session.Open()
	if err != nil {
		fmt.Printf("error opening bot: %s\n", err)
		stop <- os.Kill
	}
	defer b.session.Close()
	<-stop
	return nil
}

// After the session is open, DefaultCommands() will add
// default slash commands to the server the bot is on.
// The default commands are:
//   - /blackjack
func (b *Bot) DefaultCommands() error {
	_, err := b.session.ApplicationCommandCreate(
		b.session.State.Application.ID,
		b.session.State.Application.GuildID,
		&discordgo.ApplicationCommand{
			Name:        "blackjack",
			Description: "Launch or join a game of blackjack in the current channel",
		},
	)
	return err
}

func (b *Bot) DefaultHandlers() {
	// When the bot connects to a server
	b.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println(r.Guilds)
	})

	// When someone uses a slash command (i.e. /blackjack)
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.ApplicationCommandData().Name {
		case "blackjack":
			// check if game exists in this guild
			// if not, create a new game in this guild
			// if it exists, check if calling user is in the game
			// if not, add them to the game

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: ":white_check_mark:",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}

			ch, err := s.State.Channel(i.GuildID)
			if err != nil {
				return
			}

			if ch.IsThread() {
				return
			}

			if b.state.GuildHasGame(i.GuildID) {
				if !b.state.GameHasPlayer(i.GuildID, i.Member.User.ID) {
					thread, err := s.ThreadStart(i.ChannelID, "Blackjack Table", 0, 60)
					if err != nil {
						fmt.Println(err)
						return
					}

					header, err := s.ChannelMessageSend(thread.ID, fmt.Sprintf("<@%s>", i.Member.User.ID))
					if err != nil {
						return
					}
					b.state.AddPlayer(i.GuildID, i.Member.User.ID, thread.ID, header.ID)
				}
				return
			}

			thread, err := s.ThreadStart(i.ChannelID, "Blackjack Table", 0, 60)
			if err != nil {
				fmt.Println(err)
				return
			}

			header, err := s.ChannelMessageSend(thread.ID, fmt.Sprintf("<@%s>", i.Member.User.ID))
			if err != nil {
				return
			}

			b.state.AddNewGame(i.GuildID, thread.ID, header.ID, i.Member.User.ID)

			_, err = s.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{
				Flags:   1 << 6,
				Content: "Your move...",
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label: "Hit",
								Emoji: &discordgo.ComponentEmoji{
									Name: "💥",
								},
								Style:    discordgo.PrimaryButton,
								CustomID: "hit",
								Disabled: false,
							},
							discordgo.Button{
								Label: "Stand",
								Emoji: &discordgo.ComponentEmoji{
									Name: "🤌🏻",
								},
								Style:    discordgo.SecondaryButton,
								CustomID: "stand",
								Disabled: false,
							},
						},
					},
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	})

	// Handle Component Interactions
	b.session.AddHandler(func(s *discordgo.Session, m *discordgo.InteractionCreate) {
		switch m.Type {
		case discordgo.InteractionMessageComponent:
			if m.MessageComponentData().CustomID == "hit" {

			}

			if m.MessageComponentData().CustomID == "stand" {
			}

			err := s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: "Waiting on the table...",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label: "Hit",
									Emoji: &discordgo.ComponentEmoji{
										Name: "💥",
									},
									Style:    discordgo.PrimaryButton,
									CustomID: "hit",
									Disabled: true,
								},
								discordgo.Button{
									Label: "Stand",
									Emoji: &discordgo.ComponentEmoji{
										Name: "🤌🏻",
									},
									Style:    discordgo.SecondaryButton,
									CustomID: "stand",
									Disabled: true,
								},
							},
						},
					},
				},
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	})
}
