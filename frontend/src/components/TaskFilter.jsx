export default function TaskFilter({ filter, onFilterChange }) {
  return (
    <div className="task-filter">
      <select value={filter} onChange={(e) => onFilterChange(e.target.value)}>
        <option value="all">All Tasks</option>
        <option value="pending">Pending</option>
        <option value="in-progress">In Progress</option>
        <option value="completed">Completed</option>
      </select>
      <select onChange={(e) => onFilterChange(filter, e.target.value)}>
        <option value="all">All Priorities</option>
        <option value="low">Low</option>
        <option value="medium">Medium</option>
        <option value="high">High</option>
      </select>
    </div>
  )
}

