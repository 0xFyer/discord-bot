package bot

import (
	"fmt"
	"os"

	"github.com/0xfyer/discord-bot/game"
	"github.com/0xfyer/discord-bot/game/blackjack"
	"github.com/0xfyer/discord-bot/storage"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session *discordgo.Session
	state   *storage.State
}

func New(secret string, game game.Game) (*Bot, error) {
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
		guildID := i.GuildID
		playerID := i.Member.User.ID
		channelID := i.ChannelID

		switch i.ApplicationCommandData().Name {
		case "blackjack":
			ch, err := s.State.Channel(guildID)
			if err != nil {
				return
			}

			if ch.IsThread() {
				return
			}

			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: ":white_check_mark:",
					Flags:   1 << 6,
				},
			})
			if err != nil {
				fmt.Printf("error responsding to /blackjack: %s", err)
				return
			}

			bj := blackjack.Game{}

			// Guild has no game -- make one
			if !b.state.GuildHasGame(guildID, bj) {
				b.state.AddNewGame(guildID, bj)
			}

			// Caller is already in the game -- do nothing
			if b.state.GameHasPlayer(guildID, playerID) {
				return
			}

			// Add caller to the game and start their thread
			thread, err := s.ThreadStart(channelID, "Blackjack Table", 0, 60)
			if err != nil {
				fmt.Printf("error creating game thread: %s", err)
				return
			}

			header, err := s.ChannelMessageSend(thread.ID, fmt.Sprintf("<@%s>", playerID))
			if err != nil {
				fmt.Printf("error creating game thread header: %s", err)
				return
			}

			b.state.AddPlayer(guildID, playerID, thread.ID, header.ID)

			// Create the game thread for this player
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
				fmt.Printf("failed sending the game thread message to a new player: %s", err)
				return
			}
		}
	})

	// When someone clicks a button
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
				fmt.Printf("failed responding to a component interaction: %s", err)
				return
			}
		}

	})
}
