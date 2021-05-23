def Generate(respiration, sp02, pulse, temperature, bloodPressure, motion, proximity, overall, level, target, overallLevel):

  return dict([
    ("respiration", respiration),
    ("sp02", sp02),
    ("pulse", pulse),
    ("temperature", temperature),
    ("bloodPressure", bloodPressure),
    ("motion", motion),
    ("proximity", proximity),
  ])
