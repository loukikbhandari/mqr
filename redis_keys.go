package mqr

const (
	connectionsKey                   = "mqr::connections"                                           // Set of connection names
	connectionHeartbeatTemplate      = "mqr::connection::{connection}::heartbeat"                   // expires after {connection} died
	connectionQueuesTemplate         = "mqr::connection::{connection}::queues"                      // Set of queues consumers of {connection} are consuming
	connectionQueueConsumersTemplate = "mqr::connection::{connection}::queue::[{queue}]::consumers" // Set of all consumers from {connection} consuming from {queue}
	connectionQueueUnackedTemplate   = "mqr::connection::{connection}::queue::[{queue}]::unacked"   // List of deliveries consumers of {connection} are currently consuming

	queuesKey             = "mqr::queues"                     // Set of all open queues
	queueReadyTemplate    = "mqr::queue::[{queue}]::ready"    // List of deliveries in that {queue} (right is first and oldest, left is last and youngest)
	queueRejectedTemplate = "mqr::queue::[{queue}]::rejected" // List of rejected deliveries from that {queue}

	phConnection = "{connection}" // connection name
	phQueue      = "{queue}"      // queue name
	phConsumer   = "{consumer}"   // consumer name (consisting of tag and token)
)
