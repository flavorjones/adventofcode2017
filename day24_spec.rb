#! /usr/bin/env ruby

require 'rspec'

class Bridge
  attr_reader :components

  def initialize components
    @components = components
  end

  def strength
    @strength ||= components.inject(0) { |sum, component| sum += component.strength }
  end

  def length
    components.length
  end

  def inspect
    "[" + components.map(&:inspect).join(", ") + "]\n"
  end
end

class BridgeComponent
  attr_reader :ends

  def self.builder input_string
    input_string.split("\n").map do |line|
      BridgeComponent.new line.split("/").map(&:to_i)
    end
  end

  def initialize ends
    @ends = ends
  end

  def strength
    ends.inject(0) { |sum, plug| sum += plug }
  end

  def match? value
    ends.include? value
  end

  def other_plug value
    raise "no match" unless match?(value)
    return ends[1] if ends[0] == value
    return ends[0]
  end

  def inspect
    ends.inspect
  end
end

class BridgeBuilder
  attr_reader :components

  def initialize components
    @components = components
  end

  def bridges starting_plug
    bridges = []
    components.find_all { |component| component.match? starting_plug }.each do |component|
      bridges << Bridge.new([component])

      remaining = components - [component]
      next_starting_plug = component.other_plug(starting_plug)

      BridgeBuilder.new(remaining).bridges(next_starting_plug).each do |bridge|
        bridges << Bridge.new([component] + bridge.components)
      end
    end
    bridges
  end

  def strongest starting_plug
    bridges(starting_plug).max_by { |b| b.strength }
  end

  def longest starting_plug
    bridges(starting_plug).max_by { |b| b.length * 1000000 + b.strength }
  end
end

describe "BridgeBuilder" do
  test_input = "0/2
2/2
2/3
3/4
3/5
0/1
10/1
9/10
"
  
  let(:components) { BridgeComponent.builder(test_input) }

  it "parses components" do
    expect(components.map(&:ends)).to contain_exactly([0,2], [2,2], [2,3], [3,4], [3,5], [0,1], [10,1], [9,10])
  end

  describe "#bridges" do
    it "generates all valid bridges" do
      expect(BridgeBuilder.new(components).bridges(0).length).to eq(11)
    end
  end

  describe "#strongest" do
    it "finds the strongest bridge" do
      expect(BridgeBuilder.new(components).strongest(0).strength).to eq(31)
    end
  end

  describe "#longest" do
    it "finds the longest bridge, choosing strongest in case of tie" do
      expect(BridgeBuilder.new(components).longest(0).strength).to eq(19)
    end
  end

  describe "puzzle" do
    it "solves star 1" do
      bb = BridgeBuilder.new(BridgeComponent.builder(File.read("day24.txt")))
      strongest = bb.strongest(0)
      puts "d24 s1: strongest bridge has strength #{strongest.strength}"
    end

    it "solves star 2" do
      bb = BridgeBuilder.new(BridgeComponent.builder(File.read("day24.txt")))
      longest = bb.longest(0)
      puts "d24 s2: longest bridge has strength #{longest.strength}"
    end
  end
end
