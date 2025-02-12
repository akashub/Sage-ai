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
      <div className="space-y-3">
        {features.map((feature) => (
          <div
            key={feature.title}
            className="flex items-start gap-3 rounded-lg bg-white/[0.02] hover:bg-white/[0.04] p-4 border border-white/[0.03] transition-colors"
          >
            <div className="text-xl select-none">{feature.icon}</div>
            <div>
              <h3 className="text-[15px] font-medium text-white mb-0.5">{feature.title}</h3>
              <p className="text-[13px] leading-relaxed text-white/50">{feature.description}</p>
            </div>
          </div>
        ))}
      </div>
    )
  }
  
  