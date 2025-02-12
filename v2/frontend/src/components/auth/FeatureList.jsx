const features = [
    {
      icon: "⚡️",
      title: "Save Hours of Work",
      description: "Turn simple text prompts into ready-to-publish SQL queries in minutes, not days",
    },
    {
      icon: "📊",
      title: "Tell Better Stories",
      description: "Create engaging database queries that help you understand your data better",
    },
    {
      icon: "🎯",
      title: "Stand Out on Social",
      description: "Get more insights with professional-looking queries that capture attention",
    },
  ]
  
  export const FeatureList = () => {
    return (
      <div className="space-y-4">
        {features.map((feature) => (
          <div
            key={feature.title}
            className="flex items-start gap-4 rounded-lg bg-white/5 p-4 transition-colors hover:bg-white/10 border border-white/10"
          >
            <div className="text-xl">{feature.icon}</div>
            <div>
              <h3 className="font-medium text-white">{feature.title}</h3>
              <p className="text-sm text-white/60">{feature.description}</p>
            </div>
          </div>
        ))}
      </div>
    )
  }
  
  