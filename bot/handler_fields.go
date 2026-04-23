package bot

// This file documents the fields added to Handler to support the reminder feature.
// It extends the Handler struct defined in handler.go via a companion approach.
// Since Go does not allow splitting struct definitions across files, the reminder
// store field is declared here as documentation; the actual field must be added
// to the Handler struct in handler.go.
//
// Required addition to Handler struct in handler.go:
//
//   reminders *ReminderStore
//   prefix    string
//
// Required addition to NewHandler in handler.go:
//
//   h.reminders = NewReminderStore()
//   h.prefix = cfg.Prefix  // or "!" as default
//
// Example updated NewHandler signature (no changes to external API):
//
//   func NewHandler(cfg *config.Config, session *discordgo.Session) *Handler {
//       h := &Handler{
//           config:    cfg,
//           session:   session,
//           reminders: NewReminderStore(),
//           prefix:    cfg.Prefix,
//       }
//       return h
//   }

// ReminderHandlerVersion tracks the reminder feature version for diagnostics.
const ReminderHandlerVersion = "1.0.0"
